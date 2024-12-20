package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// MemoryStorage структура хранилища в памяти.
type MemoryStorage struct {
	// ссылки (ключ короткая ссылка, значение полная)
	db    *map[string]models.URL
	users map[int]models.User
	// удалённый url (ключ короткая ссылка, значение время)
	deletedURLs map[string]time.Time
	// ссылки пользователя (ключ короткая ссылка, значение - uuid пользователя)
	userURLs map[string]string
	// Синхронизация конккуретного доступа
	mx            sync.RWMutex
	lastIDForURL  uint
	lastIDForUser int
}

// NewMemoryStorage конструктор хранилища.
func NewMemoryStorage() *MemoryStorage {
	databaseData := make(map[string]models.URL, 1000)
	// Демо данные
	databaseData["e98192e19505472476a49f10388428ab"] = models.URL{
		ID:       1,
		ShortURL: "e98192e19505472476a49f10388428ab",
		URL:      "https://ya.ru",
	}

	instance := MemoryStorage{
		db:          &databaseData,
		users:       make(map[int]models.User, 100),
		deletedURLs: make(map[string]time.Time, 100),
		userURLs:    make(map[string]string, 100),
	}

	return &instance
}

// Add добавление нового значения.
func (s *MemoryStorage) Add(url models.URL) (int64, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	data := *s.db
	if _, ok := data[url.ShortURL]; ok {
		duplicateKeyError := pgconn.PgError{
			Code: CodeErrorDuplicateKey,
		}
		err := errors.New("short url already exists")
		return 0, errors.Join(err, &duplicateKeyError)
	}
	s.lastIDForURL++
	url.ID = s.lastIDForURL
	data[url.ShortURL] = url
	return int64(url.ID), nil
}

// CreateUser создает пользователя.
func (s *MemoryStorage) CreateUser(user models.User) (int64, error) {
	s.lastIDForUser++
	user.ID = s.lastIDForUser
	s.users[user.ID] = user
	return int64(user.ID), nil

}

// SoftDeletedShortURL Отметка об удалении ссылки.
func (s *MemoryStorage) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	for _, value := range shortURL {
		s.deletedURLs[value] = time.Now()
	}
	return nil
}

// LikeURLToUser Связывание URL с пользователем.
func (s *MemoryStorage) LikeURLToUser(urlID int64, userUUID string) error {
	for shortURL, value := range *s.db {
		if int64(value.ID) == urlID {
			s.userURLs[shortURL] = userUUID
		}
	}
	return nil
}

// MultiAdd Вставка массива.
func (s *MemoryStorage) MultiAdd(urls []models.URL) error {
	for _, url := range urls {
		s.removeItemByURL(url.URL)
		_, _ = s.Add(url)
	}
	return nil
}

// FindByShortURL поиск по короткой ссылке.
func (s *MemoryStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	data := *s.db
	if url, ok := data[shortURL]; ok {
		if deletedTime, ok2 := s.deletedURLs[shortURL]; ok2 {
			url.DeletedAt = deletedTime
		}
		return &url, nil
	}

	return nil, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL.
func (s *MemoryStorage) FindByURL(url string) (*models.URL, error) {
	var urlModel models.URL
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, modelURL := range *s.db {
		if modelURL.URL == url {
			return &modelURL, nil
		}
	}
	return &urlModel, fmt.Errorf("the url link was not found")
}

func (s *MemoryStorage) removeItemByURL(url string) {
	for shortURL, modelURL := range *s.db {
		if modelURL.URL == url {
			delete(*s.db, shortURL)
		}
	}
}

// Ping проверка доступности.
func (s *MemoryStorage) Ping() error {
	return nil
}

// FindUserByLoginAndPasswordHash Поиск пользователя.
func (s *MemoryStorage) FindUserByLoginAndPasswordHash(login string, password string) (*models.User, error) {
	for _, value := range s.users {
		if value.Login == login && value.Password == password {
			return &value, nil
		}
	}
	return nil, nil
}

// FindUrlsByUserID поиск URL-s.
func (s *MemoryStorage) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	urls := make([]models.URL, 0, 100)
	for shortURL, uuid := range s.userURLs {
		if uuid != userUUID {
			continue
		}
		if url, ok := (*s.db)[shortURL]; ok {
			urls = append(urls, url)
		}
	}
	return &urls, nil
}

// GetCountShortURL кол-во сокращенных URL
func (s *MemoryStorage) GetCountShortURL() (int64, error) {
	return int64(len(*(s.db))), nil
}

// GetCountUser кол-во пользвателей
func (s *MemoryStorage) GetCountUser() (int64, error) {
	return int64(len(s.users)), nil
}
