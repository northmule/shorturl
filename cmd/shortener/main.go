package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/certificate"
	"github.com/northmule/shorturl/internal/app/services/certificate/signers"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// @Title Shortener API
// @Description Сервис сокращения URL
// @Version 1.0
// @host      localhost:8080
func main() {
	printLabel()
	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer appStop()
	if err := run(appCtx); err != nil {
		log.Fatal(err)
	}
}

// run преднастройка
func run(ctx context.Context) error {
	err := logger.InitLogger("info")
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	storage, err := appStorage.NewStorage(ctx, cfg)
	if err != nil {
		return err
	}
	sessionStorage := appStorage.NewSessionStorage()
	shortURLService := url.NewShortURLService(storage, storage)
	stop := make(chan struct{})
	worker := workers.NewWorker(storage, stop)

	// Собираем роутер
	handlerBuilder := handlers.GetBuilder()
	handlerBuilder.SetService(shortURLService)
	handlerBuilder.SetStorage(storage)
	handlerBuilder.SetSessionStorage(sessionStorage)
	handlerBuilder.SetWorker(worker)
	handlerBuilder.SetFinderStats(storage)
	handlerBuilder.SetConfigApp(cfg)
	routes := handlerBuilder.GetAppRoutes().Init()

	if cfg.PprofEnabled {
		routes.Mount("/debug", middleware.Profiler())
	}

	httpServer := http.Server{
		Addr:    cfg.ServerURL,
		Handler: routes,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	go func() {
		<-ctx.Done()
		// Отправка сигнала о завершении в канал воркерам
		stop <- struct{}{}
		logger.LogSugar.Info("Получин сигнал. Останавливаю сервер...")

		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			logger.LogSugar.Error(err)
		}
	}()

	if cfg.EnableHTTPS {

		httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
		logger.LogSugar.Info("Подготова сертификата и ключа для TLS сервера")
		certService := certificate.NewCertificate(signers.NewEcdsaSigner())
		err = certService.InitSelfSigned()
		if err != nil {
			return err
		}
		logger.LogSugar.Infof("Сертификат: %s, ключ: %s созданы", certService.CertPath(), certService.KeyPath())
		logger.LogSugar.Infof("Running server TLS on - %s", cfg.ServerURL)
		err = httpServer.ListenAndServeTLS(certService.CertPath(), certService.KeyPath())
	} else {
		logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
		err = httpServer.ListenAndServe()
	}

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.LogSugar.Info("Сервер остановлен")
			return nil
		}
		return err
	}

	return nil
}

func printLabel() {
	template := `
	Build version: <buildVersion>
	Build date: <buildDate>
	Build commit: <buildCommit>
`
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	template = strings.Replace(template, "<buildVersion>", buildVersion, 1)
	template = strings.Replace(template, "<buildDate>", buildDate, 1)
	template = strings.Replace(template, "<buildCommit>", buildCommit, 1)

	fmt.Println(template)
}
