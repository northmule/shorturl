package handlers

import (
	"bytes"
	"github.com/northmule/shorturl/cmd/client"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestShortenerHandler тест обработчика для декодирования ссылки
func TestShortenerHandler(t *testing.T) {
	shortURLService := url.NewShortURLService(storage.NewStorage())
	ts := httptest.NewServer(AppRoutes(&shortURLService))
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
			name:    "#1_пустой_запрос_возвращает_status_bad_request",
			request: request{},
			want: want{
				code:    http.StatusBadRequest,
				isError: true,
			},
		},
		{
			name: "#2_в_body_переданна_не_ссылка_возвращается_status_bad_request",
			request: request{
				body: "Жил был слон!",
			},
			want: want{
				code:    http.StatusBadRequest,
				isError: true,
			},
		},
		{
			name: "#3_короткая_ссылка_создаётся",
			request: request{
				body: "https://ya.ru",
			},
			want: want{
				code:    http.StatusCreated,
				isError: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewBufferString(tt.request.body))
			if err != nil {
				t.Error(err)
			}
			request.Header.Set("Content-Type", "text/plain")

			response, err := client.ClientApp(client.Params{Request: request})
			if err != nil {
				t.Error(err)
			}
			defer response.Body.Close()

			respBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error(err)
			}
			if respBody == nil {
				t.Error("Тело запроса не должно быть пустым")
			}
			if tt.want.code != response.StatusCode {
				t.Errorf("Не верный код ответа сервера. Ожидается %#v пришло %#v", tt.want.code, response.StatusCode)
			}
			urlModel, _ := shortURLService.Storage.FindByURL(tt.request.body)

			if tt.want.isError == (urlModel != nil) {
				t.Error("URL не найден")
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	shortURLService := url.NewShortURLService(storage.NewStorage())
	ts := httptest.NewServer(AppRoutes(&shortURLService))
	defer ts.Close()

	request, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := ts.Client().Do(request)
	if err != nil {
		t.Error(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Ожидается код ответа StatusMethodNotAllowed")
	}
}
