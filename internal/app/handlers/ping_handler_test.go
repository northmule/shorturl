package handlers

import (
	"errors"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type MockPostgresStorageOk struct {
	mock.Mock
}

func (m *MockPostgresStorageOk) Add(url models.URL) error {
	return nil
}
func (m *MockPostgresStorageOk) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageOk) FindByURL(url string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageOk) Ping() error {
	return nil
}
func (m *MockPostgresStorageOk) MultiAdd(url []models.URL) error {
	return nil
}

type MockPostgresStorageBad struct {
	mock.Mock
}

func (m *MockPostgresStorageBad) Add(url models.URL) error {
	return nil
}
func (m *MockPostgresStorageBad) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageBad) FindByURL(url string) (*models.URL, error) {
	return nil, nil
}
func (m *MockPostgresStorageBad) Ping() error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockPostgresStorageBad) MultiAdd(url []models.URL) error {
	return nil
}

func TestPingHandler_CheckStorageConnect(t *testing.T) {
	_ = logger.NewLogger("fatal")

	file, err := os.CreateTemp("/tmp", "TestFileStorage_Add_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	postgresStorage := new(MockPostgresStorageOk)
	postgresStorage.On("Ping").Return("Ok")

	tests := []struct {
		name     string
		storage  url.StorageInterface
		wantBody string
		wantCode int
	}{
		{
			name:     "MemoryStorage",
			storage:  storage.NewMemoryStorage(),
			wantBody: "Ok",
			wantCode: http.StatusOK,
		},
		{
			name:     "MemoryStorage",
			storage:  storage.NewFileStorage(file),
			wantBody: "Ok",
			wantCode: http.StatusOK,
		},
		{
			name:     "PostgresStorage",
			storage:  postgresStorage,
			wantBody: "Ok",
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortURLService := url.NewShortURLService(tt.storage)
			ts := httptest.NewServer(AppRoutes(shortURLService))
			defer ts.Close()

			request, err := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)

			if err != nil {
				t.Fatal(err)
			}

			response, err := ts.Client().Do(request)
			if err != nil {
				t.Fatal(err)
			}
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, response.StatusCode)
			}
			if string(responseBody) != tt.wantBody {
				t.Errorf("want %s; got %s", tt.wantBody, string(responseBody))
			}
			response.Body.Close()
		})
	}

	t.Run("Возврат_ошибки_подключения", func(t *testing.T) {
		mockStorage := new(MockPostgresStorageBad)
		mockStorage.On("Ping").Return(errors.New("bad test request"))

		shortURLService := url.NewShortURLService(mockStorage)
		ts := httptest.NewServer(AppRoutes(shortURLService))
		defer ts.Close()

		request, err := http.NewRequest(http.MethodGet, ts.URL+"/ping", nil)

		if err != nil {
			t.Fatal(err)
		}

		response, err := ts.Client().Do(request)
		if err != nil {
			t.Fatal(err)
		}
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != http.StatusInternalServerError {
			t.Errorf("want %d; got %d", http.StatusInternalServerError, response.StatusCode)
		}
		if strings.Contains(string(responseBody), "bad test request") {
			t.Errorf("want %s; got %s", "bad test request", string(responseBody))
		}
		response.Body.Close()
	})
}