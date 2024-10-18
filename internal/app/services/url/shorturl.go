package url

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// ShortURLDefaultSize размер короткой ссылки.
const ShortURLDefaultSize = 10

// ShortURLData параметры данных сервиса.
type ShortURLData struct {
	URL       string
	ShortURL  string
	URLID     int64
	DeletedAt time.Time
}

// ShortURLService сервис сокращения ссылок.
type ShortURLService struct {
	Storage      IStorage
	shortURLData ShortURLData
}

// IStorage методы.
type IStorage interface {
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
}

// NewShortURLService конструктор сервиса.
func NewShortURLService(storage IStorage) *ShortURLService {
	service := &ShortURLService{
		Storage: storage,
	}

	return service
}

// DecodeURL вернёт короткий url.
func (s *ShortURLService) DecodeURL(url string) (data *ShortURLData, err error) {
	modelURL, _ := s.Storage.FindByURL(url)
	if modelURL.ShortURL != "" {
		s.shortURLData.ShortURL = modelURL.ShortURL
	} else {
		s.shortURLData.ShortURL = newRandomString(ShortURLDefaultSize)
	}

	s.shortURLData.URL = url
	urlID, err := s.Storage.Add(models.URL{
		ShortURL: s.shortURLData.ShortURL,
		URL:      s.shortURLData.URL,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != storage.CodeErrorDuplicateKey {
			logger.LogSugar.Errorf("не удалось сохранить URL %s", url)
		}
		return nil, err
	}
	s.shortURLData.URLID = urlID
	return &s.shortURLData, nil
}

// DecodeURLs преобразование массива url.
func (s *ShortURLService) DecodeURLs(urls []string) ([]models.URL, error) {
	modelURLs := make([]models.URL, len(urls))
	modelURL := new(models.URL)
	for i, url := range urls {
		modelURL.URL = url
		modelURL.ShortURL = newRandomString(ShortURLDefaultSize)
		modelURLs[i] = *modelURL
	}
	err := s.Storage.MultiAdd(modelURLs)
	if err != nil {
		return nil, err
	}
	return modelURLs, nil
}

// EncodeShortURL вернёт полный url.
func (s *ShortURLService) EncodeShortURL(shortURL string) (data *ShortURLData, err error) {
	modelURL, err := s.Storage.FindByShortURL(shortURL)
	if err != nil {
		return nil, errors.New("short url not found")
	}
	s.shortURLData.URL = modelURL.URL
	s.shortURLData.DeletedAt = modelURL.DeletedAt
	return &s.shortURLData, nil
}

func newRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz09")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
