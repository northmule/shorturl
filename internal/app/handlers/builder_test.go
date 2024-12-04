package handlers

import (
	"testing"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

func TestRoutesBuilder_SetService(t *testing.T) {
	builder := NewRoutesBuilder()
	service := &url.ShortURLService{}
	builder.SetService(service)
	if builder.shortURLService != service {
		t.Errorf("Expected shortURLService to be %v, but got %v", service, builder.shortURLService)
	}
}

func TestRoutesBuilder_SetSessionStorage(t *testing.T) {
	builder := NewRoutesBuilder()
	sessionStore := storage.NewSessionStorage()
	builder.SetSessionStorage(sessionStore)
	if builder.sessionStorage != sessionStore {
		t.Errorf("Expected sessionStorage to be %v, but got %v", sessionStore, builder.sessionStorage)
	}
}

func TestRoutesBuilder_SetWorker(t *testing.T) {
	builder := NewRoutesBuilder()
	worker := &workers.Worker{}
	builder.SetWorker(worker)
	if builder.worker != worker {
		t.Errorf("Expected worker to be %v, but got %v", worker, builder.worker)
	}
}

func TestRoutesBuilder_SetStorage(t *testing.T) {
	builder := NewRoutesBuilder()
	store := storage.NewMemoryStorage()
	builder.SetStorage(store)
	if builder.storage != store {
		t.Errorf("Expected storage to be %v, but got %v", store, builder.storage)
	}
}

func TestRoutesBuilder_SetFinderStats(t *testing.T) {
	builder := NewRoutesBuilder()
	store := storage.NewMemoryStorage()
	builder.SetFinderStats(store)
	if builder.finderStats != store {
		t.Errorf("Expected finderStats to be %v, but got %v", store, builder.finderStats)
	}
}

func TestRoutesBuilder_SetConfigApp(t *testing.T) {
	builder := NewRoutesBuilder()
	cfg := new(config.Config)
	builder.SetConfigApp(cfg)
	if builder.configApp != cfg {
		t.Errorf("Expected configApp to be %v, but got %v", cfg, builder.configApp)
	}
}

func TestRoutesBuilder_GetAppRoutes(t *testing.T) {
	logger.InitLogger("fatal")
	cfg := new(config.Config)
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	sessionStore := storage.NewSessionStorage()

	store := storage.NewMemoryStorage()
	service := url.NewShortURLService(store, store)
	worker := workers.NewWorker(store, stop)
	builder := NewRoutesBuilder()
	builder.SetService(service)
	builder.SetSessionStorage(sessionStore)
	builder.SetWorker(worker)
	builder.SetStorage(store)
	builder.SetFinderStats(store)
	builder.SetConfigApp(cfg)

	routes := builder.GetAppRoutes()

	if routes.shortURLService == nil {
		t.Errorf("Expected shortURLService to be %v, but got %v", service, routes.shortURLService)
	}
	if routes.sessionStorage == nil {
		t.Errorf("Expected sessionStorage to be %v, but got %v", sessionStore, routes.sessionStorage)
	}
	if routes.worker == nil {
		t.Errorf("Expected worker to be %v, but got %v", worker, routes.worker)
	}
	if routes.storage == nil {
		t.Errorf("Expected storage to be %v, but got %v", store, routes.storage)
	}
	if routes.finderStats == nil {
		t.Errorf("Expected finderStats to be %v, but got %v", store, routes.finderStats)
	}
	if routes.configApp == nil {
		t.Errorf("Expected configApp to be %v, but got %v", cfg, routes.configApp)
	}
}
