package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFinder struct {
	mock.Mock
}

func (m *MockFinder) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	args := m.Called(userUUID)
	return args.Get(0).(*[]models.URL), args.Error(1)
}

type MockFinderBad struct {
	mock.Mock
}

func (m *MockFinderBad) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, errors.New("error")
}

type MockDeleter struct {
	mock.Mock
	IsDeleted bool
}

func (w *MockDeleter) Del(userUUID string, input []string) {
	w.IsDeleted = true
}
func TestView(t *testing.T) {
	_ = logger.InitLogger("fatal")
	mockFinder := new(MockFinder)
	userUUID := "user123"
	userURLs := &[]models.URL{
		{ShortURL: "short1", URL: "http://example.com"},
		{ShortURL: "short2", URL: "http://example.org"},
	}
	handler := &UserURLsHandler{
		finder: mockFinder,
	}

	mockFinder.On("FindUrlsByUserID", userUUID).Return(userURLs, nil)

	req := httptest.NewRequest("GET", "/view", nil)
	ctx := context.WithValue(req.Context(), AppContext.KeyContext, userUUID)
	res := httptest.NewRecorder()
	req = req.WithContext(ctx)
	handler.View(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var responseList []ResponseView
	err := json.Unmarshal(res.Body.Bytes(), &responseList)
	assert.NoError(t, err)

	expectedResponse := []ResponseView{
		{ShortURL: "/short1", OriginalURL: "http://example.com"},
		{ShortURL: "/short2", OriginalURL: "http://example.org"},
	}
	assert.Equal(t, expectedResponse, responseList)

	mockFinder.AssertExpectations(t)
}

func BenchmarkView(b *testing.B) {
	_ = logger.InitLogger("fatal")
	mockFinder := new(MockFinder)
	userUUID := "user123"
	userURLs := &[]models.URL{
		{ShortURL: "short1", URL: "http://example.com/1"},
		{ShortURL: "short2", URL: "http://example.org/2"},
		{ShortURL: "short3", URL: "http://example.org/3"},
		{ShortURL: "short4", URL: "http://example.org/4"},
		{ShortURL: "short5", URL: "http://example.org/5"},
		{ShortURL: "short6", URL: "http://example.org/6"},
		{ShortURL: "short7", URL: "http://example.org/7"},
		{ShortURL: "short8", URL: "http://example.org/8"},
	}
	handler := &UserURLsHandler{
		finder: mockFinder,
	}
	mockFinder.On("FindUrlsByUserID", userUUID).Return(userURLs, nil)
	req := httptest.NewRequest("GET", "/view", nil)
	ctx := context.WithValue(req.Context(), AppContext.KeyContext, userUUID)
	res := httptest.NewRecorder()
	req = req.WithContext(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.View(res, req)
	}
}

func TestUserURLsHandler_Delete(t *testing.T) {
	_ = logger.InitLogger("fatal")
	deleter := new(MockDeleter)
	handler := &UserURLsHandler{
		worker: deleter,
	}

	requestBody := RequestDelete{"short1", "short2"}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	userUUID := "user123"
	req = req.WithContext(context.WithValue(req.Context(), AppContext.KeyContext, userUUID))

	handler.Delete(res, req)

	if res.Code != http.StatusAccepted {
		t.Errorf("Expected status code %d, but got %d", http.StatusAccepted, res.Code)
	}

	if !deleter.IsDeleted {
		t.Error("Deleter should be marked as deleted")
	}
}

func TestDelete_BadBody(t *testing.T) {
	_ = logger.InitLogger("fatal")
	deleter := new(MockDeleter)
	handler := &UserURLsHandler{
		worker: deleter,
	}

	errBody := io.NopCloser(&errorReader{})
	req, err := http.NewRequest("DELETE", "/api/user/urls", errBody)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, res.Code)
	}

}

func TestDelete_BadJson(t *testing.T) {
	_ = logger.InitLogger("fatal")
	deleter := new(MockDeleter)
	handler := &UserURLsHandler{
		worker: deleter,
	}

	req, err := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, res.Code)
	}

}

func TestView_BadFinder(t *testing.T) {
	_ = logger.InitLogger("fatal")
	finder := new(MockFinderBad)
	handler := &UserURLsHandler{
		finder: finder,
	}

	req, err := http.NewRequest("GET", "/api/user/urls", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler.View(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, res.Code)
	}

}

func TestView_StatusNoContent(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	handler := &UserURLsHandler{
		finder: memoryStorage,
	}

	req, err := http.NewRequest("GET", "/api/user/urls", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err)
	}
	userUUID := "user123"
	req = req.WithContext(context.WithValue(req.Context(), AppContext.KeyContext, userUUID))
	res := httptest.NewRecorder()

	handler.View(res, req)

	if res.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, but got %d", http.StatusNoContent, res.Code)
	}
}
