package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/services"
	"net/http"
)

// EncodeHandler обработчик получения оригинальной ссылки из короткой
func EncodeHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
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
