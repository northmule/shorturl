package handlers

import (
	"bytes"
	"fmt"
	"github.com/northmule/shorturl/configs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
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
		code     int
		response string
	}
	type request struct {
		method      string
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test #1 - негативный",
			request: request{
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "method not expect\n",
			},
		},
		{
			name: "Test #2 - негативный",
			request: request{
				method: http.MethodPost,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "expected Content-Type: text/plain\n",
			},
		},
		{
			name: "Test #3 - негативный",
			request: request{
				method:      http.MethodPost,
				contentType: "text/plain",
				body:        "Жил был слон!",
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "expected url\n",
			},
		},
		{
			name: "Test #4 - позитивный",
			request: request{
				method:      http.MethodPost,
				contentType: "text/plain",
				body:        "https://ya.ru",
			},
			want: want{
				code:     http.StatusCreated,
				response: fmt.Sprintf("%s/%s", configs.ServerURL, "e98192e19505472476a49f10388428ab"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.request.method, ts.URL, bytes.NewBufferString(tt.request.body))
			require.NoError(t, err)
			request.Header.Set("Content-Type", tt.request.contentType)
			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			respBody, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			assert.NotNil(t, respBody)
			assert.Equal(t, tt.want.code, response.StatusCode, "Не верный код ответа сервера")
			assert.Equal(t, tt.want.response, string(respBody), "Ошибка в значение body")
		})
	}
}
