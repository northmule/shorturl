package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
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

func TestView(t *testing.T) {
	_ = logger.NewLogger("fatal")
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
	_ = logger.NewLogger("fatal")
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
