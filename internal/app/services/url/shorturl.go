package url

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"math/rand"
	"time"
)

const ShortURLDefaultSize = 10

type ShortURLData struct {
	URL       string
	ShortURL  string
	URLID     int64
	DeletedAt time.Time
}

type ShortURLService struct {
	Storage      StorageInterface
	shortURLData ShortURLData
}

// StorageInterface методы
type StorageInterface interface {
	Add(url models.URL) (int64, error)
	CreateUser(user models.User) (int64, error)
	LikeURLToUser(urlID int64, userUUID string) error
	FindByShortURL(shortURL string) (*models.URL, error)
	FindByURL(url string) (*models.URL, error)
	Ping() error
	MultiAdd(urls []models.URL) error
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
}

func NewShortURLService(storage StorageInterface) *ShortURLService {
	service := &ShortURLService{
		Storage: storage,
	}

	return service
}

// DecodeURL вернёт короткий url
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

// DecodeURLs преобразование массива url
func (s *ShortURLService) DecodeURLs(urls []string) ([]models.URL, error) {
	modelURLs := make([]models.URL, 0)
	for _, url := range urls {
		modelURLs = append(modelURLs, models.URL{
			URL:      url,
			ShortURL: newRandomString(ShortURLDefaultSize),
		})
	}
	err := s.Storage.MultiAdd(modelURLs)
	if err != nil {
		return nil, err
	}
	return modelURLs, nil
}

// EncodeShortURL вернёт полный url
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
