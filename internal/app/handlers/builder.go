package handlers

import (
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

// RoutesBuilder простроитель объекта.
type RoutesBuilder struct {
	shortURLService *url.ShortURLService
	sessionStorage  storage.SessionAdapter
	worker          *workers.Worker
	storage         storage.StorageQuery
}

// Builder строитель.
type Builder interface {
	SetService(*url.ShortURLService)
	SetSessionStorage(storage.SessionAdapter)
	SetWorker(*workers.Worker)
	SetStorage(storage.StorageQuery)
	GetAppRoutes() *Routes
}

// NewRoutesBuilder конструктор.
func NewRoutesBuilder() *RoutesBuilder {
	return &RoutesBuilder{}
}

// GetBuilder новый объект.
func GetBuilder() Builder {
	return NewRoutesBuilder()
}

// SetService добавить сервис.
func (r *RoutesBuilder) SetService(service *url.ShortURLService) {
	r.shortURLService = service
}

// SetSessionStorage хранилище сессий.
func (r *RoutesBuilder) SetSessionStorage(adapter storage.SessionAdapter) {
	r.sessionStorage = adapter
}

// SetWorker воркер задач.
func (r *RoutesBuilder) SetWorker(worker *workers.Worker) {
	r.worker = worker
}

// SetStorage база данных.
func (r *RoutesBuilder) SetStorage(query storage.StorageQuery) {
	r.storage = query
}

// GetAppRoutes собранный роутер приложения.
func (r *RoutesBuilder) GetAppRoutes() *Routes {
	return &Routes{
		shortURLService: r.shortURLService,
		sessionStorage:  r.sessionStorage,
		worker:          r.worker,
		storage:         r.storage,
	}
}
