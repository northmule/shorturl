package handlers

import (
	"fmt"
	"github.com/northmule/shorturl/configs"
	"github.com/northmule/shorturl/internal/app/services"
	"io"
	"net/http"
	"regexp"
)

// DecodeHandler обработчик создания короткой ссылки
func DecodeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "expected post request", http.StatusBadRequest)
		return
	}
	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "expected Content-Type: text/plain", http.StatusBadRequest)
		return
	}

	postBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read body", http.StatusBadRequest)
		return
	}

	urlRegex := regexp.MustCompile(`(http|https)://\S+`)

	// Проверяем, содержится ли в postBody URL
	if !urlRegex.MatchString(string(postBody)) {
		http.Error(res, "expected url", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	shortUrlService := services.ShortUrlService{
		Url: string(postBody),
	}
	shortUrlValue, err := shortUrlService.DecodeUrl()
	if err != nil {
		http.Error(res, "error decode url", http.StatusBadRequest)
		return
	}
	shortUrl := fmt.Sprintf("%s/%s", configs.ServerUrl, shortUrlValue)
	_, err = res.Write([]byte(shortUrl))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
