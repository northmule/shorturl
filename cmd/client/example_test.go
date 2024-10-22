package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
)

var shortURLService = url.NewShortURLService(storage.NewMemoryStorage())
var stop = make(chan struct{})

func Example() {
	_ = logger.NewLogger("fatal")
	memoryStorage := storage.NewMemoryStorage()
	sessionStorage := storage.NewSessionStorage()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts := httptest.NewServer(handlers.NewRoutes(shortURLService, memoryStorage, sessionStorage).Init(ctx, stop))

	request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten/batch", bytes.NewBufferString(`[{"correlation_id":"1","original_url":"http://ya.ru"},{"correlation_id":"2","original_url":"http://ya.ru/2"},{"correlation_id":"3","original_url":"http://ya.ru/3"},{"correlation_id":"4","original_url":"http://ya.ru/4"}]`))
	if err != nil {
		logger.LogSugar.Error(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	res, err := ClientApp(Params{Request: request})
	res.Body.Close()
	if err != nil {
		logger.LogSugar.Error(err)
	}
	fmt.Print("ok")
	// Output:
	// ok
}
