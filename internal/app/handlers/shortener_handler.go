package handlers

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/services/url"
	"io"
	"net/http"
	"regexp"
)

type ShortenerHandler struct {
	regexpURLMustCompile *regexp.Regexp
	service              ShortURLServiceInterface
}

type ShortenerHandlerInterface interface {
	ShortenerHandler(res http.ResponseWriter, req *http.Request)
}

type ShortURLServiceInterface interface {
	DecodeURL(url string) (*url.ShortURLData, error)
	EncodeShortURL(string) (*url.ShortURLData, error)
}

func NewShortenerHandler(urlService ShortURLServiceInterface) ShortenerHandler {
	shortenerHandler := &ShortenerHandler{
		regexpURLMustCompile: regexp.MustCompile(`(http|https)://\S+`),
		service:              urlService,
	}
	return *shortenerHandler
}

// ShortenerHandler обработчик создания короткой ссылки
func (s *ShortenerHandler) ShortenerHandler(res http.ResponseWriter, req *http.Request) {

	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Проверяем, содержится ли в bodyValue URL
	if !s.regexpURLMustCompile.Match(bodyValue) {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")

	shortURLData, err := s.service.DecodeURL(string(bodyValue))
	if err != nil {
		http.Error(res, "error decode url", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
	shortURL := fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURLData.ShortURL)
	_, err = res.Write([]byte(shortURL))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
