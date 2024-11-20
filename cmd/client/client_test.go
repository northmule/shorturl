// Package client_test тестирует функцию ClientApp.
package client_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/northmule/shorturl/cmd/client"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/workers"
)

func TestClientApp(t *testing.T) {
	t.Run("with_request", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "https://ya.ru", nil)
		params := client.Params{Request: request}
		response, err := client.ClientApp(params)
		if err != nil {
			t.Errorf("ClientApp returned an error: %v", err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
		}
	})

}

func TestClientApp_RequestNil(t *testing.T) {
	t.Run("request_nil", func(t *testing.T) {
		_ = logger.InitLogger("fatal")
		userInput := "https://ya.ru/map1\n"
		funcDefer, err := mockStdin(t, userInput)
		if err != nil {
			t.Fatal(err)
		}

		defer funcDefer()

		params := client.Params{}

		memoryStorage := storage.NewMemoryStorage()
		sessionStorage := storage.NewSessionStorage()
		shortURLService := url.NewShortURLService(memoryStorage, memoryStorage)
		stop := make(chan struct{})
		defer func() {
			stop <- struct{}{}
		}()

		worker := workers.NewWorker(memoryStorage, stop)

		handlerBuilder := handlers.GetBuilder()
		handlerBuilder.SetService(shortURLService)
		handlerBuilder.SetStorage(memoryStorage)
		handlerBuilder.SetSessionStorage(sessionStorage)
		handlerBuilder.SetWorker(worker)
		routes := handlerBuilder.GetAppRoutes().Init()

		httpServer := http.Server{
			Addr:    ":8080",
			Handler: routes,
		}

		go func() {
			err = httpServer.ListenAndServe()
			if err != nil {
				stop <- struct{}{}
				_ = fmt.Errorf("errors: %s", err)
			}
		}()

		go func() {
			<-stop
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = httpServer.Shutdown(shutdownCtx)
		}()
		response, err := client.ClientApp(params)
		if err != nil {
			t.Fatalf("ClientApp returned an error: %v", err)
		}
		if response != nil {
			defer response.Body.Close()
		}

		modelURL, _ := memoryStorage.FindByURL("https://ya.ru/map1")
		if modelURL == nil {
			t.Error("Expected modelURL")
		}
		if modelURL != nil && modelURL.ShortURL == "" {
			t.Error("Expected modelURL.ShortUR")
		}
	})

}

func mockStdin(t *testing.T, dummyInput string) (funcDefer func(), err error) {
	t.Helper()

	oldOsStdin := os.Stdin

	tmpfile, err := os.CreateTemp(t.TempDir(), "tmp_name_test")
	if err != nil {
		return nil, err
	}

	content := []byte(dummyInput)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	os.Stdin = tmpfile

	return func() {
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}
