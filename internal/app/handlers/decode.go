package handlers

import (
	"fmt"
	"github.com/northmule/shorturl/configs"
	"io"
	"net/http"
	"regexp"
)

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
	// todo вызов сервиса сокращения сслыки
	shortUrl := fmt.Sprintf("%s/%s", configs.ServerUrl, postBody)
	res.Write([]byte(shortUrl))
}
