package handlers

import (
	"github.com/northmule/shorturl/configs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEncodeHandler тест обработчика для декодирования ссылки
func TestEncodeHandler(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	type request struct {
		method string
		id     string
	}
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Test #1 - негативный",
			request: request{
				method: http.MethodPost,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "expected get request\n",
			},
		},
		{
			name: "Test #2 - негативный",
			request: request{
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "expected id value\n",
			},
		},
		{
			name: "Test #3 - позитивный",
			request: request{
				id:     "e98192e19505472476a49f10388428ab",
				method: http.MethodGet,
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "https://ya.ru",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, configs.ServerUrl, nil)
			request.SetPathValue("id", tt.request.id)
			response := httptest.NewRecorder()

			EncodeHandler(response, request)
			assert.Equal(t, tt.want.code, response.Code, "Не верный код ответа сервера")
			assert.Equal(t, tt.want.response, response.Body.String(), "Ошибка в значение body")
		})
	}
}
