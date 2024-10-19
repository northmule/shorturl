package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
)

// @Title Shortener API
// @Description Сервис сокращения URL
// @Version 1.0
// @host      localhost:8080
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run преднастройка
func run() error {
	err := logger.NewLogger("info")
	if err != nil {
		return err
	}
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	storage, err := getStorage(cfg)
	if err != nil {
		return err
	}

	shortURLService := url.NewShortURLService(storage)
	logger.LogSugar.Infof("Running server on - %s", cfg.ServerURL)
	stop := make(chan struct{})
	routes := handlers.AppRoutes(shortURLService, stop)

	if cfg.PprofEnabled {
		routes.Mount("/debug", middleware.Profiler())
	}

	return http.ListenAndServe(cfg.ServerURL, routes)
}

func getStorage(cfg *config.Config) (url.IStorage, error) {

	if cfg.DataBaseDsn != "" {
		s, err := appStorage.NewPostgresStorage(cfg.DataBaseDsn)
		if err != nil {
			logger.LogSugar.Errorf("Failed NewPostgresStorage dsn: %s, %s", cfg.DataBaseDsn, err)
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
