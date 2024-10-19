package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/app/workers"
)

// UserURLsHandler хэндлер отображения ссылок пользователя.
type UserURLsHandler struct {
	finder  IFinderURLs
	session *storage.Session
	worker  *workers.Worker
}

// NewUserUrlsHandler Конструктор.
func NewUserUrlsHandler(finder IFinderURLs, sessionStorage *storage.Session, worker *workers.Worker) *UserURLsHandler {
	instance := UserURLsHandler{
		finder:  finder,
		session: sessionStorage,
		worker:  worker,
	}
	return &instance
}

// ResponseView структура ответа для просмотра.
type ResponseView struct {
	ShortURL    string `json:"short_Url"`
	OriginalURL string `json:"original_url"`
}

// IFinderURLs Поиск URL-s по пользователю.
type IFinderURLs interface {
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
}

// View коротки ссылки пользователя.
// @Summary Просмотр коротких ссылок пользователя
// @Failure 500
// @Failure 400
// @Success 200 {object} ResponseView
// @Router /api/user/urls [get]
func (u *UserURLsHandler) View(res http.ResponseWriter, req *http.Request) {
	userUUID := u.getUserUUID(res, req)
	logger.LogSugar.Infof("Получен запрос на просмотр URL для пользователя с uuid: %s", userUUID)
	userURLs, err := u.finder.FindUrlsByUserID(userUUID)
	if err != nil {
		http.Error(res, "Ошибка получения ссылок пользователя", http.StatusInternalServerError)
		logger.LogSugar.Error(err)
		return
	}
	res.Header().Set("content-type", "application/json")

	if len(*userURLs) == 0 {
		logger.LogSugar.Infof("Не нашёл сокращённых ссылок для пользователя с uuid: %s", userUUID)
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var responseList []ResponseView
	for _, urlItem := range *userURLs {
		responseList = append(responseList, ResponseView{
			ShortURL:    fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, urlItem.ShortURL),
			OriginalURL: urlItem.URL,
		})
	}
	responseURLs, err := json.Marshal(responseList)
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseURLs)
	if err != nil {
		http.Error(res, "error write data", http.StatusInternalServerError)
		return
	}
}

// RequestDelete запрос на удаление адресов.
type RequestDelete []string

// Delete удаление ссылок текущего пользователя.
// @Summary Удаление ссылок пользователем
// @Failure 400
// @Success 202
// @Param Delete body RequestDelete true "объект с сылками для удаления"
// @Router /api/user/urls [delete]
func (u *UserURLsHandler) Delete(res http.ResponseWriter, req *http.Request) {
	userUUID := u.getUserUUID(res, req)
	logger.LogSugar.Infof("Получен запрос на удаление для пользователя с uuid: %s", userUUID)

	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var requestShortURLs RequestDelete
	if err = json.Unmarshal(bodyValue, &requestShortURLs); err != nil {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		return
	}

	u.worker.Del(userUUID, requestShortURLs)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusAccepted)
}

func (u *UserURLsHandler) getUserUUID(res http.ResponseWriter, req *http.Request) string {
	userIDAny := req.Context().Value(context.KeyContext)
	var userUUID string
	if id, ok := userIDAny.(string); ok {
		userUUID = id
	}
	return userUUID
}
