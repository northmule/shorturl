package storage

import (
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"sync"
)

// MemoryStorage структура хранилища
type MemoryStorage struct {
	db *map[string]models.URL
	// Синхронизация конккуретного доступа
	mx sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	databaseData := make(map[string]models.URL, 1000)
	// Демо данные
	databaseData["e98192e19505472476a49f10388428ab"] = models.URL{
		ID:       1,
		ShortURL: "e98192e19505472476a49f10388428ab",
		URL:      "https://ya.ru",
	}

	instance := MemoryStorage{
		db: &databaseData,
	}

	return &instance
}

// Add добавление нового значения
func (s *MemoryStorage) Add(url models.URL) (int64, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	data := *s.db
	data[url.ShortURL] = url
	return 1, nil
}

func (s *MemoryStorage) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (s *MemoryStorage) SoftDeletedShortURL(shortURL string) error {
	return nil
}

func (s *MemoryStorage) LikeURLToUser(urlID int64, userUUID string) error {
	//todo
	return nil
}

func (s *MemoryStorage) MultiAdd(urls []models.URL) error {
	for _, url := range urls {
		s.removeItemByURL(url.URL)
		_, err := s.Add(url)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (s *MemoryStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	data := *s.db
	if url, ok := data[shortURL]; ok {
		return &url, nil
	}

	return nil, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL
func (s *MemoryStorage) FindByURL(url string) (*models.URL, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, modelURL := range *s.db {
		if modelURL.URL == url {
			return &modelURL, nil
		}
	}
	return nil, fmt.Errorf("the url link was not found")
}

func (s *MemoryStorage) removeItemByURL(url string) {
	for shortURL, modelURL := range *s.db {
		if modelURL.URL == url {
			delete(*s.db, shortURL)
		}
	}
}

func (s *MemoryStorage) GetAll() (*map[string]models.URL, error) {
	return s.db, nil
}

func (s *MemoryStorage) Ping() error {
	return nil
}

func (s *MemoryStorage) FindUserByLoginAndPasswordHash(login string, password string) (*models.User, error) {
	return nil, nil
}
func (s *MemoryStorage) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}
