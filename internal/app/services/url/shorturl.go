package url

import (
	"errors"
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"math/rand"
	"time"
)

const ShortURLDefaultSize = 6

type ShortURLData struct {
	URL      string
	ShortURL string
}

type ShortURLService struct {
	Storage      repositoryURLInterface
	shortURLData ShortURLData
}

// RepositoryURLInterface методы
type repositoryURLInterface interface {
	Add(url models.URL) error
	FindByShortURL(shortURL string) (*models.URL, error)
	FindByURL(url string) (*models.URL, error)
}

func NewShortURLService(storage repositoryURLInterface) *ShortURLService {
	service := &ShortURLService{
		Storage: storage,
	}

	return service
}

// DecodeURL вернёт короткий url
func (s *ShortURLService) DecodeURL(url string) (data *ShortURLData, err error) {
	s.shortURLData.ShortURL = newRandomString(ShortURLDefaultSize)
	s.shortURLData.URL = url
	err = s.Storage.Add(models.URL{
		ShortURL: s.shortURLData.ShortURL,
		URL:      s.shortURLData.URL,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось сохранить URL %s", url)
	}
	return &s.shortURLData, nil
}

// EncodeShortURL вернёт полный url
func (s *ShortURLService) EncodeShortURL(shortURL string) (data *ShortURLData, err error) {
	modelURL, err := s.Storage.FindByShortURL(shortURL)
	if err != nil {
		return nil, errors.New("short url not found")
	}
	s.shortURLData.URL = modelURL.URL
	return &s.shortURLData, nil
}

func newRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz09_-.")

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}

	return string(b)
}
