package main

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/filestorage"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"log"
	"net/http"
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
	config.Init()
	storage := appStorage.NewStorage()
	restoreStorageData(config.AppConfig.FileStoragePath, storage)

	var shortURLService handlers.ShortURLServiceInterface

	if config.AppConfig.FileStoragePath != "" {
		shortURLService, err = filestorage.NewSetter(config.AppConfig.FileStoragePath, url.NewShortURLService(storage))
		if err != nil {
			return err
		}
	} else {
		shortURLService = url.NewShortURLService(storage)
	}

	fmt.Println("Running server on - ", config.AppConfig.ServerURL)
	return http.ListenAndServe(config.AppConfig.ServerURL, handlers.AppRoutes(shortURLService))
}

// restoreStorageData загрузка данных URL из файла
func restoreStorageData(file string, storage *appStorage.Storage) {
	if file == "" {
		return
	}
	fileStorage, err := filestorage.NewGetter(file)
	if err != nil {
		log.Fatal(err)
	}
	storageData, err := fileStorage.ReadURLAll()
	if err != nil {
		logger.Log.Sugar().Info(err)
		return
	}
	if storageData == nil {
		logger.Log.Sugar().Info("storageData empty")
		return
	}
	storage.RestoreDBStorage(storageData)
}
