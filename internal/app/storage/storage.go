package storage

import (
	"context"
	"os"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/db"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// Storage Общие интерфейс всех методов хранилищ
type Storage interface {
	// Add добавляет URL.
	Add(url models.URL) (int64, error)
	// CreateUser создание пользователя.
	CreateUser(user models.User) (int64, error)
	// LikeURLToUser связывает пользователя с ссылкой.
	LikeURLToUser(urlID int64, userUUID string) error
	// FindByShortURL поиск по короткой ссылке.
	FindByShortURL(shortURL string) (*models.URL, error)
	// FindByURL поиск по URL.
	FindByURL(url string) (*models.URL, error)
	// Ping проверка соединения с БД.
	Ping() error
	// MultiAdd вставка массива адресов.
	MultiAdd(urls []models.URL) error
	// FindUrlsByUserID поиск ссылок пользователя
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
	// SoftDeletedShortURL пометка ссылки как удалённой.
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
	// GetCountShortURL количество коротких ссылок
	GetCountShortURL() (int64, error)
	// GetCountUser количество пользователей
	GetCountUser() (int64, error)
}

// NewStorage Создаёт нужный storage
func NewStorage(ctx context.Context, cfg *config.Config) (Storage, error) {
	if cfg.DataBaseDsn != "" {
		s, err := NewPostgresStorage(cfg.DataBaseDsn)
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
		s := NewFileStorage(file)
		return s, nil
	}

	s := NewMemoryStorage()
	return s, nil
}
