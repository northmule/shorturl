package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

// Routes маршруты приложения.
type Routes struct {
	shortURLService *url.ShortURLService
	storage         url.IStorage
	sessionStorage  storage.Session
}

// NewRoutes Конструктор маршрутов.
func NewRoutes(shortURLService *url.ShortURLService, storage url.IStorage, sessionStorage storage.Session) *Routes {
	return &Routes{
		shortURLService: shortURLService,
		storage:         storage,
		sessionStorage:  sessionStorage,
	}
}

// Init создаёт маршруты.
func (routes *Routes) Init(ctx context.Context, stop <-chan struct{}) chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})

	checkAuth := middlewarehandler.NewCheckAuth(routes.storage, routes.sessionStorage)

	r.Use(middlewarehandler.MiddlewareLogger)
	r.Use(middlewarehandler.MiddlewareGzipCompressor)

	shortenerHandler := NewShortenerHandler(routes.shortURLService, routes.storage)
	redirectHandler := NewRedirectHandler(routes.shortURLService)
	pingHandler := NewPingHandler(routes.storage)

	worker := workers.NewWorker(routes.storage, stop)
	userUrlsHandler := NewUserUrlsHandler(routes.storage, routes.sessionStorage, worker)

	r.With(
		checkAuth.AuthEveryone,
	).Post("/", shortenerHandler.ShortenerHandler)
	r.Get("/{id}", redirectHandler.RedirectHandler)
	r.With(
		checkAuth.AuthEveryone,
	).Post("/api/shorten", shortenerHandler.ShortenerJSONHandler)
	r.Get("/ping", pingHandler.CheckStorageConnect)
	r.Post("/api/shorten/batch", shortenerHandler.ShortenerBatch)

	r.With(
		checkAuth.AccessVerificationUserUrls,
		checkAuth.AuthEveryone,
	).Get("/api/user/urls", userUrlsHandler.View)

	r.With(
		checkAuth.AccessVerificationUserUrls,
		checkAuth.AuthEveryone,
	).Delete("/api/user/urls", userUrlsHandler.Delete)

	return r
}
