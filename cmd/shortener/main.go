package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/db"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

// @Title Shortener API
// @Description Сервис сокращения URL
// @Version 1.0
// @host      localhost:8080
func main() {
	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
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

	storage, err := getStorage(ctx, cfg)
	if err != nil {
		return err
	}
	sessionStorage := appStorage.NewSessionStorage()
	shortURLService := url.NewShortURLService(storage, storage)
	stop := make(chan struct{})
	worker := workers.NewWorker(storage, stop)
	routes := handlers.NewRoutes(shortURLService, storage, sessionStorage, worker).Init()

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

	logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
	err = httpServer.ListenAndServe()

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.LogSugar.Info("Сервер остановлен")
			return nil
		}
		return err
	}

	return nil
}

func getStorage(ctx context.Context, cfg *config.Config) (appStorage.StorageQuery, error) {

	if cfg.DataBaseDsn != "" {
		s, err := appStorage.NewPostgresStorage(cfg.DataBaseDsn)
		if err != nil {
			logger.LogSugar.Errorf("Failed NewPostgresStorage dsn: %s, %s", cfg.DataBaseDsn, err)
			return nil, err
		}

		logger.LogSugar.Info("Инициализация миграций")
		migrations := db.NewMigrations(s.RawDB)
		err = migrations.Up(ctx)
		if err != nil {
			return nil, err
		}

		return s, nil
	}

	if cfg.FileStoragePath != "" {
		file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.LogSugar.Errorf("Failed to open file %s: error: %s", cfg.FileStoragePath, err)
			return nil, err
		}
		return appStorage.NewFileStorage(file), nil
	}

	return appStorage.NewMemoryStorage(), nil
}
