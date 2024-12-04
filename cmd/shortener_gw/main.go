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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/certificate"
	"github.com/northmule/shorturl/internal/app/services/certificate/signers"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	grpcHandlers "github.com/northmule/shorturl/internal/grpc/handlers"
	"github.com/northmule/shorturl/internal/grpc/handlers/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const (
	gRPCGatewayServerAddress = "localhost:8081"
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

	logger.LogSugar.Info("создаём gRPC-сервер")
	authInterceptor := interceptors.NewCheckAuth(storage, sessionStorage)
	trustedInterceptor := interceptors.NewCheckTrustedSubnet(cfg)
	loggerInterceptor := interceptors.NewLogger(logger.LogSugar)

	mux := runtime.NewServeMux()

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor([]grpc.UnaryServerInterceptor{
		loggerInterceptor.LogStart,
		authInterceptor.AuthEveryone,
		authInterceptor.AccessVerificationUserUrls,
		trustedInterceptor.GrantAccess,
		loggerInterceptor.LogEnd,
	}...))

	logger.LogSugar.Info("Подготовка сервисов")
	contract.RegisterPingHandlerServer(grpcServer, grpcHandlers.NewPingHandler(storage))
	contract.RegisterRedirectHandlerServer(grpcServer, grpcHandlers.NewRedirectHandler(shortURLService))
	contract.RegisterShortenerHandlerServer(grpcServer, grpcHandlers.NewShortenerHandler(shortURLService, storage, storage))
	contract.RegisterStatsHandlerServer(grpcServer, grpcHandlers.NewStatsHandler(storage))
	contract.RegisterUserUrlsHandlerServer(grpcServer, grpcHandlers.NewUserURLsHandler(storage, sessionStorage, worker))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = errors.Join(contract.RegisterPingHandlerHandlerFromEndpoint(ctx, mux, gRPCGatewayServerAddress, opts))
	err = errors.Join(err, contract.RegisterRedirectHandlerHandlerFromEndpoint(ctx, mux, gRPCGatewayServerAddress, opts))
	err = errors.Join(err, contract.RegisterShortenerHandlerHandlerFromEndpoint(ctx, mux, gRPCGatewayServerAddress, opts))
	err = errors.Join(err, contract.RegisterStatsHandlerHandlerFromEndpoint(ctx, mux, gRPCGatewayServerAddress, opts))
	err = errors.Join(err, contract.RegisterUserUrlsHandlerHandlerFromEndpoint(ctx, mux, gRPCGatewayServerAddress, opts))

	if err != nil {
		return err
	}

	grpcLc := net.ListenConfig{}
	grpcListen, err := grpcLc.Listen(ctx, "tcp", gRPCGatewayServerAddress)
	if err != nil {
		return err
	}

	go func() {
		logger.LogSugar.Infof("Running gRPC server on - %s", gRPCGatewayServerAddress)
		if err = grpcServer.Serve(grpcListen); err != nil {
			return
		}
	}()

	httpServer := http.Server{
		Addr:    cfg.ServerURL,
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		// Отправка сигнала о завершении в канал воркерам
		stop <- struct{}{}
		logger.LogSugar.Info("Получин сигнал. Останавливаю HTTP сервер...")

		shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()

		logger.LogSugar.Info("Получин сигнал. Останавливаю gRPC сервер...")
		grpcServer.GracefulStop()

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
		logger.LogSugar.Infof("Running HTTP server TLS on - %s", cfg.ServerURL)
		err = httpServer.ListenAndServeTLS(certService.CertPath(), certService.KeyPath())
	} else {
		logger.LogSugar.Infof("Running HTTP server on - %s", cfg.ServerURL)
		err = httpServer.ListenAndServe()
	}

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.LogSugar.Info("Сервер HTTP остановлен")
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
