package storage

import (
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"sync"
)

// Storage структура хранилища
type Storage struct {
	db *map[string]models.URL
	// Синхронизация конккуретного доступа
	mx sync.RWMutex
}

func NewStorage() *Storage {
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

// RestoreDBStorage восстановит бд из переданного значения
func (s *Storage) RestoreDBStorage(db map[string]models.URL) {
	s.db = &db
}
