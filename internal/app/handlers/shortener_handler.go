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

var regexURL = regexp.MustCompile(`(http|https)://\S+`)

type ShortenerHandler struct {
	regexURL *regexp.Regexp
	service  *url.ShortURLService
}

type ShortenerHandlerInterface interface {
	ShortenerHandler(res http.ResponseWriter, req *http.Request)
}

func NewShortenerHandler(urlService *url.ShortURLService) ShortenerHandler {
	shortenerHandler := &ShortenerHandler{
		regexURL: regexURL,
		service:  urlService,
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
	if !s.regexURL.Match(bodyValue) {
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

type ShortenerRequest struct {
	URL string `json:"URL"`
}
type JSONResponse struct {
	Result string `json:"result"`
}

// ShortenerJSONHandler принимает и отдаёт json
func (s *ShortenerHandler) ShortenerJSONHandler(res http.ResponseWriter, req *http.Request) {

	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var shortenerRequest ShortenerRequest
	if err = json.Unmarshal(bodyValue, &shortenerRequest); err != nil {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		return
	}

	// Проверяем, содержится ли в bodyValue URL
	if !s.regexURL.MatchString(shortenerRequest.URL) {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "application/json")

	shortURLData, err := s.service.DecodeURL(shortenerRequest.URL)
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

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenerBatch обработка списка адресов
func (s *ShortenerHandler) ShortenerBatch(res http.ResponseWriter, req *http.Request) {

	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var requestItems []BatchRequest
	if err = json.Unmarshal(bodyValue, &requestItems); err != nil {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		return
	}
	urls := make([]string, 0)
	for _, requestItem := range requestItems {
		if !s.regexURL.MatchString(requestItem.OriginalURL) {
			continue
		}
		urls = append(urls, requestItem.OriginalURL)
	}
	modelURLs, err := s.service.DecodeURLs(urls)
	if err != nil {
		http.Error(res, "error decode urls", http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "application/json")

	responseItems := make([]BatchResponse, 0)
	for _, requestItem := range requestItems {
		for _, modelURL := range modelURLs {
			if requestItem.OriginalURL == modelURL.URL {
				responseItems = append(responseItems, BatchResponse{
					CorrelationID: requestItem.CorrelationID,
					ShortURL:      fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, modelURL.ShortURL),
				})
			}
		}
	}

	responseString, err := json.Marshal(responseItems)
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
