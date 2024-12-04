package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
	"github.com/stretchr/testify/assert"
)

type mockBadUserFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockBadUserFinder) GetCountShortURL() (int64, error) {
	return 1, nil
}

// GetCountUser кол-во пользвателей
func (s *mockBadUserFinder) GetCountUser() (int64, error) {
	return 0, errors.New("error")
}

type mockBadURLsFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockBadURLsFinder) GetCountShortURL() (int64, error) {
	return 0, errors.New("error")
}

// GetCountUser кол-во пользвателей
func (s *mockBadURLsFinder) GetCountUser() (int64, error) {
	return 1, nil
}

type mockFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockFinder) GetCountShortURL() (int64, error) {
	return 1, nil
}

// GetCountUser кол-во пользвателей
func (s *mockFinder) GetCountUser() (int64, error) {
	return 1, nil
}

func TestStatsHandler_ViewStats(t *testing.T) {

	memoryStorage := storage.NewMemoryStorage()

	tests := []struct {
		name         string
		finder       StatsFinder
		expectedCode int
	}{
		{
			name:         "error_GetCountUser",
			finder:       new(mockBadUserFinder),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "error_GetCountShortURL",
			finder:       new(mockBadURLsFinder),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "ok",
			finder:       new(mockFinder),
			expectedCode: http.StatusOK,
		},
	}

	_ = logger.InitLogger("fatal")

	stor := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(stor, stor)
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	ts := httptest.NewServer(NewRoutes(shortURLService, stor, storage.NewSessionStorage(), workers.NewWorker(memoryStorage, stop)).Init())
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewStatsHandler(tt.finder)
			req, err := http.NewRequest(http.MethodGet, ts.URL+"/api/internal/stats", nil)
			if err != nil {
				t.Error(err)
			}
			res := httptest.NewRecorder()
			h.ViewStats(res, req)
			assert.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
