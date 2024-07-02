package handlers

import (
	"errors"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestIteration2_EncodeHandler тест обработчика для декодирования ссылки
func TestIteration2_EncodeHandler(t *testing.T) {
	ts := httptest.NewServer(AppRoutes())
	defer ts.Close()

	type want struct {
		code     int
		location string
	}
	type request struct {
		id string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test #1 - позитивный",
			request: request{
				id: "e98192e19505472476a49f10388428ab",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://ya.ru",
			},
		},
		{
			name: "Test #2 - негативный",
			request: request{
				id: "123",
			},
			want: want{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}

	config.AppConfig.DatabasePath = "shorturl_test.db"
	err := storage.AutoMigrate()
	require.NoError(t, err)
	defer func() {
		_ = os.Remove("shorturl_test.db")
	}()

	appStorage := storage.New()

	url := models.URL{
		ShortURL: "e98192e19505472476a49f10388428ab",
		URL:      "https://ya.ru",
	}
	err = appStorage.Add(&url)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Отключить переход по ссылке при положительном ответе сервиса
			ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
				errorRedirect := errors.New("HTTP redirect blocked")
				return errorRedirect
			}
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/"+tt.request.id, nil)

			require.NoError(t, err)

			response, err := ts.Client().Do(request)
			response.Body.Close()
			var errorValue string
			if err != nil {
				errorValue = err.Error()
			}

			if errorValue != "" && strings.Contains(errorValue, "HTTP redirect blocked") {
				err = nil
			}
			require.NoError(t, err)

			location := response.Header.Get("Location")
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode, "Не верный код ответа сервера")
			assert.Equal(t, tt.want.location, location, "Ошибка в значение body")
		})
	}
}
