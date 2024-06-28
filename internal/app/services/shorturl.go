package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage"
)

type ShortURLService struct {
	URL      string
	ShortURL string
}

// ShortURLService вернёт короткий url
func (s *ShortURLService) DecodeURL() (shortURL string, err error) {
	algData := md5.Sum([]byte(s.URL)) // todo заменить
	return fmt.Sprintf("%x", algData), nil
}

// EncodeShortURL вернёт полный url
func (s *ShortURLService) EncodeShortURL() (shortURL string, err error) {
	modelURL, ok := storage.DatabaseData[s.ShortURL]
	if !ok {
		return "", errors.New("short url not found")
	}
	return modelURL.URL, nil
}
