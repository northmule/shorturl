package handlers

import (
	"github.com/go-chi/chi/v5"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"net/http"
)

// EncodeHandler обработчик получения оригинальной ссылки из короткой
func EncodeHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(res, "expected id value", http.StatusBadRequest)
		return
	}
	appStorage := appStorage.New()
	modelURL, err := appStorage.FindByShortURL(id)
	if err != nil {
		http.Error(res, "error encode shortUrl", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.Header().Set("Location", modelURL.URL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
