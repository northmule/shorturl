package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func ExamplePingHandler_CheckStorageConnect() {
	url := "http://localhost:8080/ping"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "gophermart_session=20106c00221e4fe38582e36a12327fdfc05b904bbad8d30d238a8f8323fcf90d:fbbad27c-16b3-48e3-a455-785074e45981; shorturl_session=2e41a78f9851029e40e85e70ac38d24ca06dcc8bcbb1da0e3619f17edd6c050a:431300ac-c58c-4dcf-941c-e47ca43511ba")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func ExampleRedirectHandler_RedirectHandler() {
	url := "http://localhost:8080/hRBwJoehFV1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "gophermart_session=20106c00221e4fe38582e36a12327fdfc05b904bbad8d30d238a8f8323fcf90d:fbbad27c-16b3-48e3-a455-785074e45981; shorturl_session=2e41a78f9851029e40e85e70ac38d24ca06dcc8bcbb1da0e3619f17edd6c050a:431300ac-c58c-4dcf-941c-e47ca43511ba")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func ExampleShortenerHandler_ShortenerHandler() {
	url := "http://localhost:8080"
	method := "POST"

	payload := strings.NewReader(`https://habr.com/133`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Cookie", "gophermart_session=20106c00221e4fe38582e36a12327fdfc05b904bbad8d30d238a8f8323fcf90d:fbbad27c-16b3-48e3-a455-785074e45981; shorturl_session=2e41a78f9851029e40e85e70ac38d24ca06dcc8bcbb1da0e3619f17edd6c050a:431300ac-c58c-4dcf-941c-e47ca43511ba")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Header.Get("Content-Type"))
}

func ExampleShortenerHandler_ShortenerJSONHandler() {
	url := "http://localhost:8080/api/shorten"
	method := "POST"

	payload := strings.NewReader(`{
    "url": "https://ya.ru
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "2848ab491d6d81b8e89053213d14b50aa67993905ded671aa7280b022a070182")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "shorturl_session=2848ab491d6d81b8e89053213d14b50aa67993905ded671aa7280b022a070182:cd0f4b6d-7cef-4f56-830f-a3ed1c460d7b")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func ExampleShortenerHandler_ShortenerBatch() {

	url := "http://localhost:8080/api/shorten/batch"
	method := "POST"

	payload := strings.NewReader(`[
    {
        "correlation_id": "1",
        "original_url": "http://ya.ru/212"
    }
]`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "shorturl_session=2848ab491d6d81b8e89053213d14b50aa67993905ded671aa7280b022a070182:cd0f4b6d-7cef-4f56-830f-a3ed1c460d7b")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func ExampleUserURLsHandler_View() {
	url := "http://localhost:8080/api/user/urls"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "shorturl_session=2848ab491d6d81b8e89053213d14b50aa67993905ded671aa7280b022a070182:cd0f4b6d-7cef-4f56-830f-a3ed1c460d7b")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func ExampleUserURLsHandler_Delete() {
	url := "http://localhost:8080/api/user/urls"
	method := "DELETE"

	payload := strings.NewReader(`["BsJPu9gh9N"]`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "shorturl_session=2848ab491d6d81b8e89053213d14b50aa67993905ded671aa7280b022a070182:cd0f4b6d-7cef-4f56-830f-a3ed1c460d7b")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
