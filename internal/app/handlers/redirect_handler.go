package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/services/url"
	"net/http"
)

type RedirectHandler struct {
	service *url.ShortURLService
}

type RedirectHandlerInterface interface {
	RedirectHandler(res http.ResponseWriter, req *http.Request)
}

func NewRedirectHandler(urlService *url.ShortURLService) RedirectHandler {
	redirectHandler := &RedirectHandler{
		service: urlService,
	}
	return *redirectHandler
}

// RedirectHandler обработчик получения оригинальной ссылки из короткой
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
	res.Header().Set("Location", modelURL.URL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
