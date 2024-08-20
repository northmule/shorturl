package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers/auth"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"net/http"
)

type UserURLsHandler struct {
	storage url.StorageInterface
	session *storage.SessionStorage
}

func NewUserUrlsHandler(storage url.StorageInterface, sessionStorage *storage.SessionStorage) *UserURLsHandler {
	instance := UserURLsHandler{
		storage: storage,
		session: sessionStorage,
	}
	return &instance
}

type ResponseView struct {
	ShortUrl    string `json:"short_Url"`
	OriginalUrl string `json:"original_url"`
}

// View коротки ссылки пользователя
func (u *UserURLsHandler) View(res http.ResponseWriter, req *http.Request) {
	token := auntificator.GetUserToken(req)
	if token == "" {
		res.WriteHeader(http.StatusUnauthorized)
		logger.LogSugar.Infof("Ожидалось значение cookie %s", auth.CookieAuthName)
		return
	}
	userUUID := "a4a45d8d-cd8b-47a7-a7a1-4bafcf3d83a5"

	for k, v := range u.session.Values {
		if k == token {
			userUUID = v
			break
		}
	}

	userURLs, err := u.storage.FindUrlsByUserId(userUUID)
	if err != nil {
		http.Error(res, "Ошибка получения ссылок пользователя", http.StatusInternalServerError)
		logger.LogSugar.Error(err)
		return
	}
	res.Header().Set("content-type", "application/json")

	if len(*userURLs) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var responseList []ResponseView
	for _, urlItem := range *userURLs {
		responseList = append(responseList, ResponseView{
			ShortUrl:    fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, urlItem.ShortURL),
			OriginalUrl: urlItem.URL,
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
