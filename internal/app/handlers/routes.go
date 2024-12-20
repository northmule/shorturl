package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

// Routes маршруты приложения.
type Routes struct {
	shortURLService *url.ShortURLService
	sessionStorage  storage.SessionAdapter
	worker          *workers.Worker
	storage         storage.StorageQuery
	finderStats     StatsFinder
	configApp       *config.Config
}

// todo: поменять на RoutesBuilder
// NewRoutes Конструктор маршрутов.
func NewRoutes(shortURLService *url.ShortURLService, storage storage.StorageQuery, sessionStorage storage.SessionAdapter, worker *workers.Worker) *Routes {
	return &Routes{
		shortURLService: shortURLService,
		sessionStorage:  sessionStorage,
		worker:          worker,
		storage:         storage,
	}
}

// Init создаёт маршруты.
func (routes *Routes) Init() chi.Router {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not expect\n"))
	})

	checkAuth := middlewarehandler.NewCheckAuth(routes.storage, routes.sessionStorage)
	checkTrustedSubnet := middlewarehandler.NewCheckTrustedSubnet(routes.configApp)

	r.Use(middleware.RequestLogger(logger.LogSugar))
	r.Use(middlewarehandler.MiddlewareGzipCompressor)

	shortenerHandler := NewShortenerHandler(routes.shortURLService, routes.storage, routes.storage)
	redirectHandler := NewRedirectHandler(routes.shortURLService)
	pingHandler := NewPingHandler(routes.storage)

	userUrlsHandler := NewUserUrlsHandler(routes.storage, routes.sessionStorage, routes.worker)

	statsHandler := NewStatsHandler(routes.finderStats)

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

	r.With(
		checkTrustedSubnet.GrantAccess,
	).Get("/api/internal/stats", statsHandler.ViewStats)

	return r
}
