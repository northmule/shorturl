package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"net/http"
)

// AppRoutes маршруты приложения
func AppRoutes(shortURLService ShortURLServiceInterface) chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})

	r.Use(middlewarehandler.MiddlewareLogger)

	shortenerHandler := NewShortenerHandler(shortURLService)
	redirectHandler := NewRedirectHandler(shortURLService)

	r.Post("/", shortenerHandler.ShortenerHandler)
	r.Get("/{id}", redirectHandler.RedirectHandler)
	r.Post("/api/shorten", shortenerHandler.ShortenerJsonHandler)

	return r
}
