package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIteration2_EncodeHandler тест обработчика для декодирования ссылки
func TestIteration2_EncodeHandler(t *testing.T) {
	ts := httptest.NewServer(AppRoutes())
	defer ts.Close()

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
				id:     "123",
				method: http.MethodPost,
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "method not expect\n",
			},
		},
		{
			name: "Test #2 - позитивный",
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
			request, err := http.NewRequest(tt.request.method, ts.URL+"/"+tt.request.id, nil)
			require.NoError(t, err)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)
			defer response.Body.Close()

			respBody, err := io.ReadAll(response.Body)
			stringBody := string(respBody)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode, "Не верный код ответа сервера")
			assert.Equal(t, tt.want.response, stringBody, "Ошибка в значение body")
		})
	}
}
