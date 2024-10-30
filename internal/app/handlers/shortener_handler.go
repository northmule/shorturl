package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// ShortenerHandler хэндлер сокращения ссылок.
type ShortenerHandler struct {
	service *url.ShortURLService
	finder  Finder
	setter  LikeURLToUserSetter
}

// Finder поиск значений.
type Finder interface {
	// FindUrlsByUserID поиск ссылок пользователя
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
	// FindByShortURL поиск по короткой ссылке.
	FindByShortURL(shortURL string) (*models.URL, error)
	// FindByURL поиск по URL.
	FindByURL(url string) (*models.URL, error)
}

// NewShortenerHandler конструктор.
func NewShortenerHandler(urlService *url.ShortURLService, finder Finder, setter LikeURLToUserSetter) ShortenerHandler {
	shortenerHandler := &ShortenerHandler{
		service: urlService,
		finder:  finder,
		setter:  setter,
	}
	return *shortenerHandler
}

// LikeURLToUserSetter интерфейс связывания URL с пользователем.
type LikeURLToUserSetter interface {
	LikeURLToUser(urlID int64, userUUID string) error
}

// ShortenerHandler обработчик создания короткой ссылки.
// @Summary Получение короткой ссылки
// @Failure 400
// @Success 307 {string} Location "origin_url"
// @Param url body string true "оригинальная ссылка для сокращения" SchemaExample(https://ya.ru/1)
// @Router / [post]
func (s *ShortenerHandler) ShortenerHandler(res http.ResponseWriter, req *http.Request) {
	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Проверяем, содержится ли в bodyValue URL.
	bodyString := string(bodyValue)
	if !strings.Contains(bodyString, "http://") && !strings.Contains(bodyString, "https://") {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")

	var (
		headerStatus int
		shortURL     string
	)
	userIDAny := req.Context().Value(context.KeyContext)
	var userUUID string
	if id, ok := userIDAny.(string); ok {
		userUUID = id
	}
	shortURL, headerStatus, err = s.fillShortURLAndResponseStatus(userUUID, string(bodyValue))
	if err != nil {
		http.Error(res, "error find model", headerStatus)
		return
	}

	res.WriteHeader(headerStatus)
	shortURL = fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURL)
	_, err = res.Write([]byte(shortURL))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}

// ShortenerRequest запрос к методу ShortenerJSONHandler.
type ShortenerRequest struct {
	URL string `json:"URL"`
}

// JSONResponse ответ от метода ShortenerJSONHandler.
type JSONResponse struct {
	Result string `json:"result"`
}

// ShortenerJSONHandler принимает и отдаёт json.
// @Summary Получение коротких ссылок
// @Failure 400
// @Success 200 {object} JSONResponse
// @Param ShortenerJSONHandler body ShortenerRequest true "объект с сылками для сокращения"
// @Router /api/shorten [post]
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
	if !strings.Contains(shortenerRequest.URL, "http://") && !strings.Contains(shortenerRequest.URL, "https://") {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "application/json")

	var (
		responseJSON JSONResponse
		headerStatus int
		shortURL     string
	)
	userIDAny := req.Context().Value(context.KeyContext)
	var userUUID string
	if id, ok := userIDAny.(string); ok {
		userUUID = id
	}
	shortURL, headerStatus, err = s.fillShortURLAndResponseStatus(userUUID, shortenerRequest.URL)
	if err != nil {
		http.Error(res, "error find model", headerStatus)
		return
	}

	responseJSON = JSONResponse{
		Result: fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, shortURL),
	}
	responseString, err := json.Marshal(responseJSON)
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(headerStatus)

	_, err = res.Write(responseString)
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}

// BatchRequest запрос для списка адресов.
type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponse ответ для списка адресов.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenerBatch обработка списка адресов.
// @Summary Получение коротких ссылок
// @Failure 400
// @Success 200 {object} BatchResponse
// @Param ShortenerBatch body BatchRequest true "объект с сылками для сокращения"
// @Router /api/shorten/batch [post]
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
		if !strings.Contains(requestItem.OriginalURL, "http://") && !strings.Contains(requestItem.OriginalURL, "https://") {
			continue
		}
		urls = append(urls, requestItem.OriginalURL)
	}
	if len(urls) == 0 {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		return
	}
	modelURLs, err := s.service.DecodeURLs(urls)
	if err != nil {
		http.Error(res, "error decode urls", http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "application/json")

	responseItems := make([]BatchResponse, 0, len(requestItems))
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
func (s *ShortenerHandler) fillShortURLAndResponseStatus(userUUID string, url string) (string, int, error) {
	var (
		headerStatus int
		shortURL     string
		isURLExists  bool
	)
	shortURLData, err := s.service.DecodeURL(url)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.CodeErrorDuplicateKey {
			isURLExists = true
		} else {
			return "", http.StatusInternalServerError, err
		}
	}
	if isURLExists {
		modelURL, err := s.finder.FindByURL(url)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
		headerStatus = http.StatusConflict
		shortURL = modelURL.ShortURL
	} else {
		headerStatus = http.StatusCreated
		shortURL = shortURLData.ShortURL
		err = s.setter.LikeURLToUser(shortURLData.URLID, userUUID)
		if err != nil {
			return "", 0, err
		}
	}

	return shortURL, headerStatus, nil
}
