package storage

import (
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// С sqllit тесты не проходят
//// Storage структура хранилища
//type Storage struct {
//	db *gorm.DB
//}
//
//// repositoryURL методы
//type repositoryURL interface {
//	Update(id string, url models.URL) error
//	Add(url models.URL) error
//	Get(id int) (models.URL, error)
//	FindByShortURL(shortURL string) (models.URL, error)
//	FindByURL(url string) (models.URL, error)
//}
//
//func New() Storage {
//	db, err := gorm.Open(sqlite.Open(config.AppConfig.DatabasePath), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//	storage := Storage{}
//	storage.db = db
//	return storage
//}
//
//func AutoMigrate() error {
//	storage := New()
//	err := storage.db.AutoMigrate(models.URL{})
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// Get получения значения по id
//func (s *Storage) Get(id int) (models.URL, error) {
//	modelURL := models.URL{}
//	result := New().db.First(&modelURL, id)
//	if result.Error != nil {
//		return modelURL, result.Error
//	}
//	return modelURL, nil
//}
//
//// Update обновление значения по id
//func (s *Storage) Update(id string, url models.URL) error {
//	return nil
//}
//
//// Add добавление нового значения
//func (s *Storage) Add(url *models.URL) error {
//	result := New().db.Create(&url)
//
//	return result.Error
//}
//
//// FindByShortURL поиск по короткой ссылке
//func (s *Storage) FindByShortURL(shortURL string) (models.URL, error) {
//	modelURL := models.URL{}
//	result := New().db.First(&modelURL, "short_url = ?", shortURL)
//	return modelURL, result.Error
//}
//
//// FindByURL поиск по URL
//func (s *Storage) FindByURL(url string) (models.URL, error) {
//	modelURL := models.URL{}
//	result := New().db.First(&modelURL, "url = ?", url)
//
//	return modelURL, result.Error
//}

type DatabaseMap map[string]models.URL

// Storage структура хранилища
type Storage struct {
	db *DatabaseMap
}

// DatabaseData временные данные
var DatabaseData DatabaseMap

// repositoryURL методы
type repositoryURL interface {
	Add(url models.URL) error
	FindByShortURL(shortURL string) (models.URL, error)
	FindByURL(url string) (models.URL, error)
}

func Init() {
	if len(DatabaseData) > 0 {
		return
	}
	DatabaseData = make(DatabaseMap, 500)
	DatabaseData = DatabaseMap{
		"e98192e19505472476a49f10388428ab": {
			ID:       1,
			ShortURL: "e98192e19505472476a49f10388428ab",
			URL:      "https://ya.ru",
		},
	}
}

func New() *Storage {
	Init()
	storage := Storage{}
	storage.db = &DatabaseData
	return &storage
}

// Add добавление нового значения
func (s *Storage) Add(url models.URL) error {
	data := *s.db
	data[url.ShortURL] = url
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (s *Storage) FindByShortURL(shortURL string) (models.URL, error) {
	data := *s.db
	if url, ok := data[shortURL]; ok {
		return url, nil
	}

	return models.URL{}, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL
func (s *Storage) FindByURL(url string) (models.URL, error) {
	for _, modelURL := range *s.db {
		if modelURL.URL == url {
			return modelURL, nil
		}
	}
	return models.URL{}, fmt.Errorf("the url link was not found")
}
