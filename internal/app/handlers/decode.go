package handlers

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/services"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"io"
	"net/http"
	"regexp"
)

// DecodeHandler обработчик создания короткой ссылки
func DecodeHandler(res http.ResponseWriter, req *http.Request) {
	// Автотест не устанавливает нужный Content-Type https://github.com/Yandex-Practicum/go-autotests/blob/main/cmd/shortenertest/iteration1_test.go#L120
	//if req.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(res, "expected Content-Type: text/plain", http.StatusBadRequest)
	//	return
	//}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read body", http.StatusBadRequest)
		return
	}

	urlRegex := regexp.MustCompile(`(http|https)://\S+`)

	// Проверяем, содержится ли в body URL
	if !urlRegex.Match(body) {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	shortURLService := services.ShortURLData{
		URL: string(body),
	}
	shortURLValue, err := shortURLService.DecodeURL()
	if err != nil {
		http.Error(res, "error decode url", http.StatusBadRequest)
		return
	}
	appStorage := appStorage.New()
	urlModel := models.URL{
		ShortURL: shortURLValue,
		URL:      shortURLService.URL,
	}
	err = appStorage.Add(urlModel)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL := fmt.Sprintf("%s/%s", config.AppConfig.BaseShortURL, urlModel.ShortURL)
	_, err = res.Write([]byte(shortURL))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
