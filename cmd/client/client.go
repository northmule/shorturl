package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	endpoint    = "http://localhost:8080/"
	contentType = "application/x-www-form-urlencoded"
)

type Params struct {
	Request *http.Request
}

func ClientApp(params Params) (*http.Response, error) {

	request := params.Request
	isStdin := false
	if params.Request == nil {
		data := url.Values{}
		fmt.Println("Введите длинный URL")
		reader := bufio.NewReader(os.Stdin)
		long, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		long = strings.TrimSuffix(long, "\n")
		data.Set("url", long)
		request, err = http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Content-Type", contentType)
		isStdin = true
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		if isStdin {
			log.Fatal(err)
		} else {
			return nil, err
		}
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		if isStdin {
			log.Fatal(err)
		} else {
			return nil, err
		}
	}
	if isStdin {
		fmt.Println("Статус-код ", response.Status)
		fmt.Println(string(body))
		return nil, err
	}
	return response, nil
}
