package handlers

import (
	"bytes"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestIteration1Empty Для инкремента1 не должно быть тестов
func TestIteration1Empty(t *testing.T) {
	assert.True(t, true)
}

// TestIteration2_DecodeHandler тест обработчика для декодирования ссылки
func TestIteration2_DecodeHandler(t *testing.T) {
	ts := httptest.NewServer(AppRoutes())
	defer ts.Close()

	type want struct {
		code    int
		isError bool
	}
	type request struct {
		method string
		body   string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name:    "Test #1 - негативный",
			request: request{},
			want: want{
				code:    http.StatusBadRequest,
				isError: true,
			},
		},
		{
			name: "Test #2 - негативный",
			request: request{
				body: "Жил был слон!",
			},
			want: want{
				code:    http.StatusBadRequest,
				isError: true,
			},
		},
		{
			name: "Test #3 - позитивный",
			request: request{
				body: "https://ya.ru",
			},
			want: want{
				code:    http.StatusCreated,
				isError: false,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBufferString(tt.request.body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "text/plain")
			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			respBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.NotNil(t, respBody)
			assert.Equal(t, tt.want.code, response.StatusCode, "Не верный код ответа сервера")
			_, err = appStorage.FindByURL(tt.request.body)
			assert.Equal(t, tt.want.isError, err != nil)
		})
	}
}
