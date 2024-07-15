package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/services/url"
	"io"
	"net/http"
	"regexp"
)

var regexMustCompileURL = regexp.MustCompile(`(http|https)://\S+`)

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
		regexpURLMustCompile: regexMustCompileURL,
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

type JSONRequest struct {
	URL string `json:"URL"`
}
type JSONResponse struct {
	Result string `json:"result,omitempty"`
}

// ShortenerJSONHandler принимает и отдаёт json
func (s *ShortenerHandler) ShortenerJSONHandler(res http.ResponseWriter, req *http.Request) {

	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var jsonRequest JSONRequest
	if err = json.Unmarshal(bodyValue, &jsonRequest); err != nil {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		return
	}

	// Проверяем, содержится ли в bodyValue URL
	if !s.regexpURLMustCompile.MatchString(jsonRequest.URL) {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "application/json")

	shortURLData, err := s.service.DecodeURL(jsonRequest.URL)
	if err != nil {
		http.Error(res, "error decode url", http.StatusBadRequest)
		return
	}

	responseJSON := JSONResponse{
		Result: fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURLData.ShortURL),
	}
	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)

	_, err = res.Write(responseString)
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
