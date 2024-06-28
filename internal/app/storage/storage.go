package storage

import (
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// Storage Интерфейс хранения
type Storage interface {
	Get(id string) (string, error)
	Update(id string, value string) error
	Add(value string) error
	FindByShortUrl(shortUrl string) models.Url
}

// DatabaseStorage Реализация хранения в виде структуры
type DatabaseStorage struct{}

// Get получения значения по id
func (d *DatabaseStorage) Get(id string) (string, error) {
	// Реализация логики получения значения по ключу
	return "Value for " + id, nil
}

// Update обновление значения по id
func (d *DatabaseStorage) Update(id string, value string) error {
	return nil
}

// Add добавление нового значения
func (d *DatabaseStorage) Add(value string) error {
	return nil
}

// FindByShortUrl поиск по короткой ссылке
func (d *DatabaseStorage) FindByShortUrl(shortUrl string) error {
	return nil
}
