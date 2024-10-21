package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/mock"
)

type MockPostgresStorageOk struct {
	mock.Mock
}

func (m *MockPostgresStorageOk) Add(url models.URL) (int64, error) {
	return 0, nil
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
func (m *MockPostgresStorageOk) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (m *MockPostgresStorageOk) LikeURLToUser(urlID int64, userUUID string) error {
	return nil
}

func (m *MockPostgresStorageOk) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (m *MockPostgresStorageOk) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

type MockPostgresStorageBad struct {
	mock.Mock
}

func (m *MockPostgresStorageBad) Add(url models.URL) (int64, error) {
	return 0, nil
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
func (m *MockPostgresStorageBad) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (m *MockPostgresStorageBad) LikeURLToUser(urlID int64, userUUID string) error {
	return nil
}

func (m *MockPostgresStorageBad) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (m *MockPostgresStorageBad) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
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
		storage  url.IStorage
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
			name:     "FileStorage",
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
			stop := make(chan struct{})
			ts := httptest.NewServer(AppRoutes(shortURLService, stop))
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
		stop := make(chan struct{})
		ts := httptest.NewServer(AppRoutes(shortURLService, stop))
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
