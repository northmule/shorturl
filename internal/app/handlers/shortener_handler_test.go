package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/northmule/shorturl/cmd/client"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
)

// TestShortenerHandler тест обработчика для декодирования ссылки
func TestShortenerHandler(t *testing.T) {
	_, _ = logger.NewLogger("info")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
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

			if tt.want.isError == (urlModel.URL != "") {
				t.Error("URL не найден")
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
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

func TestShortenerJsonHandler(t *testing.T) {
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
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
				body: `{"url":"https://ya.ru"}`,
			},
			want: want{
				code:    http.StatusCreated,
				isError: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBufferString(tt.request.body))
			if err != nil {
				t.Error(err)
			}
			request.Header.Set("Content-Type", "application/json")

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
				t.Error("Тело ответа не должно быть пустым")
			}
			if tt.want.code != response.StatusCode {
				t.Errorf("Не верный код ответа сервера. Ожидается %#v пришло %#v", tt.want.code, response.StatusCode)
			}
			var jsonResponse JSONResponse
			err = json.Unmarshal(respBody, &jsonResponse)
			if err != nil && !tt.want.isError {
				t.Errorf("Ошибка разбора json ответа: %s", respBody)
			}
			jsonResponse.Result = strings.Trim(jsonResponse.Result, "/")
			urlModel, _ := shortURLService.Storage.FindByShortURL(jsonResponse.Result)

			if tt.want.isError == (urlModel != nil) {
				t.Error("URL не найден")
			}
		})
	}
}

func TestGzipCompression(t *testing.T) {
	_, _ = logger.NewLogger("info")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
	defer ts.Close()

	requestBody := `{
    		"url": "https://ya.ru"
	}`

	t.Run("Отправка_серверу_от_клиента_сжатых_данных_в_виде_json_строки", func(t *testing.T) {
		// Сжатие данных, для отправки с стороны браузера
		encodeBuffer := bytes.NewBuffer(nil)
		gzipBuffer := gzip.NewWriter(encodeBuffer)
		_, err := gzipBuffer.Write([]byte(requestBody))
		if err != nil {
			t.Error(err)
		}
		err = gzipBuffer.Close()
		if err != nil {
			t.Fatal(err)
		}
		encodeString := encodeBuffer.String()

		if encodeString == "" {
			t.Fatal("Запрос не был сжат перед отправкой")
		}
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", encodeBuffer)
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "no") // т.к клиент по умолчанию шлёт их
		response, err := client.ClientApp(client.Params{Request: request})
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		respBody, _ := io.ReadAll(response.Body)

		var jsonResponse JSONResponse
		err = json.Unmarshal(respBody, &jsonResponse)
		if err != nil {
			t.Errorf("Ошибка разбора json ответа: %s", respBody)
		}
		jsonResponse.Result = strings.Trim(jsonResponse.Result, "/")
		// Если всё ок, то должена найтись модель по короткой ссылке с сервера
		urlModel, _ := shortURLService.Storage.FindByShortURL(jsonResponse.Result)
		if urlModel == nil {
			t.Error("Закодированный URL из ответа в БД не найден")
		}
	})

	t.Run("не_отправляем_поддерживаемый_content_type, сжатия_не_должно_быть", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBufferString(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Accept-Encoding", "gzip")
		response, err := client.ClientApp(client.Params{Request: request})
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		respBody, _ := io.ReadAll(response.Body)

		var jsonResponse JSONResponse
		err = json.Unmarshal(respBody, &jsonResponse)
		if err != nil {
			t.Errorf("Ошибка разбора json ответа: %s", respBody)
		}
		jsonResponse.Result = strings.Trim(jsonResponse.Result, "/")
		// Если всё ок, то должена найтись модель по короткой ссылке с сервера
		urlModel, _ := shortURLService.Storage.FindByShortURL(jsonResponse.Result)
		if urlModel == nil {
			t.Error("Закодированный URL из ответа в БД не найден")
		}
	})

	t.Run("проверяем_что_сервер_вернул_тело_ответа_в_сжатом_виде", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBufferString(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept-Encoding", "gzip")
		response, err := client.ClientApp(client.Params{Request: request})
		if err != nil {
			t.Error(err)
		}
		defer response.Body.Close()

		unpackBody, err := gzip.NewReader(response.Body) // Распаковываем данные ответа
		if err != nil {
			t.Fatal(err)
		}
		respBody, err := io.ReadAll(unpackBody)
		if err != nil {
			t.Fatal(err)
		}
		var jsonResponse JSONResponse
		err = json.Unmarshal(respBody, &jsonResponse)
		if err != nil {
			t.Errorf("Ошибка разбора json ответа: %s", respBody)
			return
		}
		jsonResponse.Result = strings.Trim(jsonResponse.Result, "/")
		// Если всё ок, то должена найтись модель по короткой ссылке с сервера
		urlModel, _ := shortURLService.Storage.FindByShortURL(jsonResponse.Result)
		if urlModel == nil {
			t.Error("Закодированный URL из ответа в БД не найден")
		}
	})
}

func BenchmarkShortenerHandler(b *testing.B) {
	_, _ = logger.NewLogger("fatal")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
	defer ts.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/", bytes.NewBufferString("https://ya.ru"))
		if err != nil {
			b.Error(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		b.StartTimer()
		res, err := client.ClientApp(client.Params{Request: request})
		res.Body.Close()
		if err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkShortenerJSONHandler(b *testing.B) {
	_, _ = logger.NewLogger("fatal")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
	defer ts.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", bytes.NewBufferString(`{"url":"https://ya.ru"}`))
		if err != nil {
			b.Error(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		b.StartTimer()
		res, err := client.ClientApp(client.Params{Request: request})
		res.Body.Close()
		if err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkShortenerBatch(b *testing.B) {
	_, _ = logger.NewLogger("fatal")
	shortURLService := url.NewShortURLService(storage.NewMemoryStorage())
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(NewRoutes(shortURLService, storage.NewMemoryStorage(), storage.NewSessionStorage()).Init(ctx, stop))
	defer ts.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", bytes.NewBufferString(`[{"correlation_id":"1","original_url":"http://ya.ru"},{"correlation_id":"2","original_url":"http://ya.ru/2"},{"correlation_id":"3","original_url":"http://ya.ru/3"},{"correlation_id":"4","original_url":"http://ya.ru/4"}]`))
		if err != nil {
			b.Error(err)
		}
		request.Header.Set("Content-Type", "text/plain")
		b.StartTimer()
		res, err := client.ClientApp(client.Params{Request: request})
		res.Body.Close()
		if err != nil {
			b.Error(err)
		}

	}
}
