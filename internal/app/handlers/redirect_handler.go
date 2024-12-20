package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/services/url"
)

// RedirectHandler хэндлер для обработки коротких ссылок.
type RedirectHandler struct {
	service *url.ShortURLService
}

// NewRedirectHandler конструктор хэндлера.
func NewRedirectHandler(urlService *url.ShortURLService) RedirectHandler {
	redirectHandler := &RedirectHandler{
		service: urlService,
	}
	return *redirectHandler
}

// RedirectHandler обработчик получения оригинальной ссылки из короткой.
// @Summary Преобразование короткой ссылки в оригинальную с переходом по ссылке
// @Failure 410
// @Success 307 {string} Location "origin_url"
// @Router /{id} [get]
func (r *RedirectHandler) RedirectHandler(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(res, "expected id value", http.StatusBadRequest)
		return
	}
	modelURL, err := r.service.EncodeShortURL(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	res.Header().Set("content-type", "text/plain")
	if modelURL.DeletedAt.IsZero() {
		res.Header().Set("Location", modelURL.URL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusGone)
	}

}
