package handlers

import (
	"github.com/northmule/shorturl/internal/app/services"
	"net/http"
)

// EncodeHandler обработчик получения оригинальной ссылки из короткой
func EncodeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "expected get request", http.StatusBadRequest)
		return
	}
	id := req.PathValue("id")
	if id == "" {
		http.Error(res, "expected id value", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)

	shortURLService := services.ShortURLService{
		ShortURL: id,
	}
	urlValue, err := shortURLService.EncodeShortURL()
	if err != nil {
		http.Error(res, "error encode shortUrl", http.StatusBadRequest)
		return
	}

	_, err = res.Write([]byte(urlValue))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
