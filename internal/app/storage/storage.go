package storage

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/filestorage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"os"
	"sync"
)

// Storage структура хранилища
type Storage struct {
	db *map[string]models.URL
	// Синхронизация конккуретного доступа
	mx sync.RWMutex
}

func NewStorage(restore bool) *Storage {
	databaseData := make(map[string]models.URL, 1000)
	// Демо данные
	databaseData["e98192e19505472476a49f10388428ab"] = models.URL{
		ID:       1,
		ShortURL: "e98192e19505472476a49f10388428ab",
		URL:      "https://ya.ru",
	}

	storage := Storage{
		db: &databaseData,
	}
	if restore {
		storage.restoreStorage(config.AppConfig.FileStoragePath)
	}
	return &storage
}

// Add добавление нового значения
func (s *Storage) Add(url models.URL) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	data := *s.db
	data[url.ShortURL] = url
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (s *Storage) FindByShortURL(shortURL string) (*models.URL, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	data := *s.db
	if url, ok := data[shortURL]; ok {
		return &url, nil
	}

	return nil, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL
func (s *Storage) FindByURL(url string) (*models.URL, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, modelURL := range *s.db {
		if modelURL.URL == url {
			return &modelURL, nil
		}
	}
	return nil, fmt.Errorf("the url link was not found")
}

// restoreStorage восстановит бд из переданного значения
func (s *Storage) restoreStorage(filePath string) {
	if filePath == "" {
		logger.LogSugar.Error("path filePath empty")
		return
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.LogSugar.Errorf("Failed to open filePath %s: error: %s", filePath, err)
		return
	}
	fileStorage, err := filestorage.NewGetter(file)
	if err != nil {
		logger.LogSugar.Errorf("filestorage.NewGetter(%s) %s", filePath, err)
		return
	}
	storageData, err := fileStorage.ReadURLAll()
	if err != nil {
		logger.LogSugar.Error(err)
		return
	}
	if storageData == nil {
		logger.LogSugar.Info("storageData empty")
		return
	}
	s.db = &storageData
}
