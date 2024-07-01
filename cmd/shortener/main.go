package main

import (
	"fmt"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// run преднастройка
func run() error {
	configInit := config.Init()
	configInit.InitEnvConfig()
	configInit.InitFlagConfig()

	fmt.Println("Running server on - ", config.AppConfig.ServerURL)
	return http.ListenAndServe(config.AppConfig.ServerURL, handlers.AppRoutes())
}
