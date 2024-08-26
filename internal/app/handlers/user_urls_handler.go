package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/app/workers"
	"io"
	"net/http"
)

const defaultUUID = "a4a45d8d-cd8b-47a7-a7a1-4bafcf3d1111"

type UserURLsHandler struct {
	finder  FinderURLs
	session *storage.SessionStorage
	worker  *workers.Worker
}

func NewUserUrlsHandler(finder FinderURLs, sessionStorage *storage.SessionStorage, worker *workers.Worker) *UserURLsHandler {
	instance := UserURLsHandler{
		finder:  finder,
		session: sessionStorage,
		worker:  worker,
	}
	return &instance
}

type ResponseView struct {
	ShortURL    string `json:"short_Url"`
	OriginalURL string `json:"original_url"`
}

type FinderURLs interface {
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
}

// View коротки ссылки пользователя
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

type RequestDelete []string

func (u *UserURLsHandler) Delete(res http.ResponseWriter, req *http.Request) {
	userUUID := u.getUserUUID(res, req)
	logger.LogSugar.Infof("Получен запрос на удаление для пользователя в uuid: %s", userUUID)

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

	userURLs, err := u.finder.FindUrlsByUserID(userUUID)
	if err != nil {
		http.Error(res, "Ошибка получения ссылок пользователя", http.StatusInternalServerError)
		logger.LogSugar.Error(err)
		return
	}
	if len(*userURLs) == 0 {
		http.Error(res, "Пользователь ещё не создал ни одной ссылки", http.StatusNoContent)
		logger.LogSugar.Infof("При запросе FindUrlsByUserID(%s) не найдено URLs", userUUID)
		return
	}

	shortURLs := make([]string, 0)
	for _, userURL := range *userURLs {
		for _, requestShortURL := range requestShortURLs {
			if userURL.ShortURL == requestShortURL {
				shortURLs = append(shortURLs, requestShortURL)
				break
			}
		}
	}
	if len(shortURLs) > 0 {
		u.deleteShortLinks(shortURLs)
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusAccepted)
}

func (u *UserURLsHandler) deleteShortLinks(shortURLs []string) {

	u.worker.Del(shortURLs)
}

func (u *UserURLsHandler) getUserUUID(res http.ResponseWriter, req *http.Request) string {
	//token := auntificator.GetUserToken(req)
	//if token == "" {
	//	res.WriteHeader(http.StatusUnauthorized)
	//	logger.LogSugar.Infof("Ожидалось значение cookie %s", auntificator.CookieAuthName)
	//	return defaultUUID
	//}
	//userUUID := defaultUUID
	//
	//for k, v := range u.session.GetAll() {
	//	if k == token {
	//		userUUID = v
	//		break
	//	}
	//}
	userIDAny := req.Context().Value(context.KeyContext)
	var userUUID string
	if id, ok := userIDAny.(string); ok {
		userUUID = id
	}
	return userUUID
}
