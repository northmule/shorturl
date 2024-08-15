package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/handlers/auth"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"net/http"
)

// AppRoutes маршруты приложения
func AppRoutes(shortURLService *url.ShortURLService) chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})

	r.Use(middlewarehandler.MiddlewareLogger)
	r.Use(middlewarehandler.MiddlewareGzipCompressor)
	r.Use(middlewarehandler.CheckAuth)

	shortenerHandler := NewShortenerHandler(shortURLService)
	redirectHandler := NewRedirectHandler(shortURLService)
	pingHandler := NewPingHandler(shortURLService.Storage)
	jwtHandler := auth.NewJWTHandler(shortURLService.Storage)

	sessionStorage := storage.NewSessionStorage()
	hmacHandler := auth.NewHMACHandler(shortURLService.Storage, sessionStorage)
	userUrlsHandler := NewUserUrlsHandler(shortURLService.Storage, sessionStorage)

	r.Post("/", shortenerHandler.ShortenerHandler)
	r.Get("/{id}", redirectHandler.RedirectHandler)
	r.Post("/api/shorten", shortenerHandler.ShortenerJSONHandler)
	r.Get("/ping", pingHandler.CheckStorageConnect)
	r.Post("/api/shorten/batch", shortenerHandler.ShortenerBatch)

	r.Post("/api/auth_jwt", jwtHandler.Auth)
	r.Post("/api/auth_hmac", hmacHandler.Auth)
	r.Get("/api/auth_hmac_everyone", hmacHandler.AuthEveryone)

	r.Get("/api/user/urls", userUrlsHandler.View)

	return r
}
