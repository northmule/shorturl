package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/app/workers"
)

// TestRedirectHandler тест обработчика для декодирования ссылки
func TestRedirectHandler(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage(), workers.NewWorker(memoryStorage, stop)).Init())

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

func TestRedirectHandler_StatusGone(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage(), workers.NewWorker(memoryStorage, stop)).Init())
	defer ts.Close()
	// необходимые данные
	identy, _ := memoryStorage.Add(models.URL{
		ShortURL: "ttt",
		URL:      "https://ya.ru",
	})
	_, _ = memoryStorage.CreateUser(models.User{
		UUID: "1111-2222-3333",
	})
	_ = memoryStorage.LikeURLToUser(identy, "1111-2222-3333")
	_ = memoryStorage.SoftDeletedShortURL("1111-2222-3333", "ttt")

	request, err := http.NewRequest(http.MethodGet, ts.URL+"/ttt", nil)

	if err != nil {
		t.Error(err)
	}

	response, err := ts.Client().Do(request)
	if err != nil {
		t.Error(err)
	}
	defer response.Body.Close()

	if http.StatusGone != response.StatusCode {
		t.Errorf("Не верный код ответа сервера. Ожидается %#v пришло %#v", http.StatusGone, response.StatusCode)
	}

}

func TestRedirectHandler_StatusBadRequest(t *testing.T) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage(), workers.NewWorker(memoryStorage, stop)).Init())
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	h := NewRedirectHandler(shortURLService)
	h.RedirectHandler(res, req)

	if http.StatusBadRequest != res.Code {
		t.Errorf("Не верный код ответа сервера. Ожидается %#v пришло %#v", http.StatusBadRequest, res.Code)
	}
}

func BenchmarkRedirectHandler(b *testing.B) {
	_ = logger.InitLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)
	sessionStorage := storage.NewSessionStorage()
	stop := make(chan struct{})
	defer func() {
		stop <- struct{}{}
	}()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), sessionStorage, workers.NewWorker(memoryStorage, stop)).Init())
	// Отключить переход по ссылке при положительном ответе сервиса
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		errorRedirect := errors.New("HTTP redirect blocked")
		return errorRedirect
	}
	defer ts.Close()

	b.Run("короткая_ссылка_не_существует", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/e98192e19505472476a49f10388428ab", nil)
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
			response, _ := ts.Client().Do(request)
			response.Body.Close()
		}
	})

	b.Run("короткая_ссылка_существует", func(b *testing.B) {
		shortURL := "e98192e19505472476a49f10388428ab"
		memoryStorage.Add(models.URL{
			ShortURL: shortURL,
			URL:      "https://ya.ru/123",
		})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/"+shortURL, nil)
			if err != nil {
				b.Error(err)
			}
			b.StartTimer()
			response, _ := ts.Client().Do(request)
			response.Body.Close()
		}
	})
}
