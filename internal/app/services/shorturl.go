package services

import (
	"errors"
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage"
	"math/rand"
	"time"
)

const ShortURLDefaultSize = 6

type ShortURLData struct {
	URL      string
	ShortURL string
}

type ShortURLService interface {
	DecodeURL() (string, error)
	EncodeShortURL() (string, error)
}

// DecodeURL вернёт короткий url
func (s *ShortURLData) DecodeURL() (shortURL string, err error) {
	shortURL = newRandomString(ShortURLDefaultSize)
	return fmt.Sprintf("%x", shortURL), nil
}

// EncodeShortURL вернёт полный url
func (s *ShortURLData) EncodeShortURL() (shortURL string, err error) {
	modelURL, ok := storage.DatabaseData[s.ShortURL]
	if !ok {
		return "", errors.New("short url not found")
	}
	return modelURL.URL, nil
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
