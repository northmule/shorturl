package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage"
)

type ShortUrlService struct {
	Url      string
	ShortUrl string
}

// ShortUrlService вернёт короткий url
func (s *ShortUrlService) DecodeUrl() (shortUrl string, err error) {
	algData := md5.Sum([]byte(s.Url)) // todo заменить
	return fmt.Sprintf("%x", algData), nil
}

// EncodeShortUrl вернёт полный url
func (s *ShortUrlService) EncodeShortUrl() (shortUrl string, err error) {
	modelUrl, ok := storage.DatabaseData[s.ShortUrl]
	if !ok {
		return "", errors.New("short url not found")
	}
	return modelUrl.Url, nil
}
