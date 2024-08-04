package main

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"log"
	"net/http"
	"os"
)

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
	_, err = config.Init()
	if err != nil {
		return err
	}

	var storage url.StorageInterface

	if config.AppConfig.DataBaseDsn != "" {
		storage, err = appStorage.NewPostgresStorage(config.AppConfig.DataBaseDsn)
		if err != nil {
			logger.LogSugar.Errorf("Failed NewPostgresStorage dsn: %s, %s", config.AppConfig.DataBaseDsn, err)
			return err
		}
	} else if config.AppConfig.FileStoragePath != "" {
		file, err := os.OpenFile(config.AppConfig.FileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.LogSugar.Errorf("Failed to open file %s: error: %s", config.AppConfig.FileStoragePath, err)
			return err
		}
		storage = appStorage.NewFileStorage(file)
	} else {
		storage = appStorage.NewMemoryStorage()
	}
	shortURLService := url.NewShortURLService(storage)
	fmt.Println("Running server on - ", config.AppConfig.ServerURL)
	return http.ListenAndServe(config.AppConfig.ServerURL, handlers.AppRoutes(shortURLService))
}
