package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// AppRoutes маршруты приложения
func AppRoutes(shortURLService ShortURLServiceInterface) chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})

	shortenerHandler := NewShortenerHandler(shortURLService)
	redirectHandler := NewRedirectHandler(shortURLService)

	r.Post("/", shortenerHandler.ShortenerHandler)
	r.Get("/{id}", redirectHandler.RedirectHandler)

	return r
}
