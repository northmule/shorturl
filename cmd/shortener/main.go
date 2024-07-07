package main

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"log"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run преднастройка
func run() error {
	err := logger.NewLogger("info")
	if err != nil {
		return err
	}
	config.Init()

	shortURLService := url.NewShortURLService(appStorage.NewStorage())
	fmt.Println("Running server on - ", config.AppConfig.ServerURL)
	return http.ListenAndServe(config.AppConfig.ServerURL, handlers.AppRoutes(&shortURLService))
}
