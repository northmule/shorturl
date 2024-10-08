package handlers

import (
	"errors"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRedirectHandler тест обработчика для декодирования ссылки
func TestRedirectHandler(t *testing.T) {
	_ = logger.NewLogger("fatal")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ts := httptest.NewServer(AppRoutes(shortURLService, stop))

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
			name: "test_#1_короткая_ссылка_преобразуется_в_длинную",
			request: request{
				id: "e98192e19505472476a49f10388428ab",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				location: "https://ya.ru",
			},
		},
		{
			name: "Test_#2_передан_не_существующий_id, вернётся_status_bad_request",
			request: request{
				id: "123",
			},
			want: want{
				code:     http.StatusNotFound,
				location: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Отключить переход по ссылке при положительном ответе сервиса
			ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
				errorRedirect := errors.New("HTTP redirect blocked")
				return errorRedirect
			}
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/"+tt.request.id, nil)

			if err != nil {
				t.Error(err)
			}

			response, err := ts.Client().Do(request)
			response.Body.Close()

			var errorValue string
			if err != nil {
				errorValue = err.Error()
			}

			if errorValue != "" && strings.Contains(errorValue, "HTTP redirect blocked") {
				err = nil
			}
			if err != nil {
				t.Error(err)
			}

			location := response.Header.Get("Location")
			if err != nil {
				t.Error(err)
			}
			if tt.want.code != response.StatusCode {
				t.Errorf("Не верный код ответа сервера. Ожидается %#v пришло %#v", tt.want.code, response.StatusCode)
			}
			if tt.want.location != location {
				t.Errorf("Ошибка в значение body. Ожидается %#v пришло %#v", tt.want.code, response.StatusCode)
			}
		})
	}
}
