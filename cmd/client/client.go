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

// todo переделать
func ClientStart() {
	data := url.Values{}
	fmt.Println("Введите длинный URL")

	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	long = strings.TrimSuffix(long, "\n")
	data.Set("url", long)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", contentType)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fmt.Println("Статус-код ", response.Status)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
