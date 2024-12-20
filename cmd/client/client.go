// Package client client используется для запросов в рамках тестов приложения.
package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Настройки клиента.
const (
	endpoint    = "http://localhost:8080/"
	contentType = "application/x-www-form-urlencoded"
)

// Params параметры клиента.
type Params struct {
	// Request запрос
	Request *http.Request
}

// ClientApp Клиент для запросов в тестах
func ClientApp(params Params) (*http.Response, error) {

	request := params.Request
	isStdin := false
	if params.Request == nil {
		// data := url.Values{}
		log.Println("Введите длинный URL")
		reader := bufio.NewReader(os.Stdin)
		long, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		inputURL := strings.TrimSuffix(long, "\n")
		// data.Set("url", long)
		request, err = http.NewRequest(http.MethodPost, endpoint, strings.NewReader(inputURL))
		if err != nil {
			return nil, err
		}
		request.Header.Add("Content-Type", contentType)
		isStdin = true
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if isStdin {
		log.Println("Статус-код ", response.Status)
		body, err := io.ReadAll(response.Body)
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(response.Body)
		if err != nil {
			return nil, fmt.Errorf("execute request: %v", err)
		}
		log.Println(string(body))
		return nil, err
	}
	return response, nil
}
