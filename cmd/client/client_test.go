// Package client_test тестирует функцию ClientApp.
package client_test

import (
	"net/http"
	"testing"

	"github.com/northmule/shorturl/cmd/client"
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
