package storage

import (
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Storage структура хранилища
type Storage struct {
	db *gorm.DB
}

// repositoryURL методы
type repositoryURL interface {
	Update(id string, url models.URL) error
	Add(url models.URL) error
	Get(id int) (models.URL, error)
	FindByShortURL(shortURL string) (models.URL, error)
	FindByURL(url string) (models.URL, error)
}

func New() Storage {
	db, err := gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	storage := Storage{}
	storage.db = db
	return storage
}

func AutoMigrate() error {
	storage := New()
	err := storage.db.AutoMigrate(models.URL{})
	if err != nil {
		return err
	}
	return nil
}

// Get получения значения по id
func (s *Storage) Get(id int) (models.URL, error) {
	modelURL := models.URL{}
	result := New().db.First(&modelURL, id)
	if result.Error != nil {
		return modelURL, result.Error
	}
	return modelURL, nil
}

// Update обновление значения по id
func (s *Storage) Update(id string, url models.URL) error {
	return nil
}

// Add добавление нового значения
func (s *Storage) Add(url *models.URL) error {
	result := New().db.Create(&url)

	return result.Error
}

// FindByShortURL поиск по короткой ссылке
func (s *Storage) FindByShortURL(shortURL string) (models.URL, error) {
	modelURL := models.URL{}
	result := New().db.First(&modelURL, "short_url = ?", shortURL)
	return modelURL, result.Error
}

// FindByURL поиск по URL
func (s *Storage) FindByURL(url string) (models.URL, error) {
	modelURL := models.URL{}
	result := New().db.First(&modelURL, "url = ?", url)

	return modelURL, result.Error
}
