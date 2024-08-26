package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
	"net/http"
)

// AppRoutes маршруты приложения
func AppRoutes(shortURLService *url.ShortURLService) chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})
	sessionStorage := storage.NewSessionStorage()
	checkAuth := middlewarehandler.NewCheckAuth(shortURLService.Storage, sessionStorage)

	r.Use(middlewarehandler.MiddlewareLogger)
	r.Use(middlewarehandler.MiddlewareGzipCompressor)
	r.Use(checkAuth.AuthEveryone)

	shortenerHandler := NewShortenerHandler(shortURLService, shortURLService.Storage)
	redirectHandler := NewRedirectHandler(shortURLService)
	pingHandler := NewPingHandler(shortURLService.Storage)

	worker := workers.NewWorker(shortURLService.Storage)
	userUrlsHandler := NewUserUrlsHandler(shortURLService.Storage, sessionStorage, worker)

	r.Post("/", shortenerHandler.ShortenerHandler)
	r.Get("/{id}", redirectHandler.RedirectHandler)
	r.Post("/api/shorten", shortenerHandler.ShortenerJSONHandler)
	r.Get("/ping", pingHandler.CheckStorageConnect)
	r.Post("/api/shorten/batch", shortenerHandler.ShortenerBatch)

	r.Get("/api/user/urls", userUrlsHandler.View)
	r.Delete("/api/user/urls", userUrlsHandler.Delete)

	return r
}
