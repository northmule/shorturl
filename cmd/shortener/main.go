package main

import (
	"github.com/northmule/shorturl/configs"
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
	return http.ListenAndServe(configs.ServerURL, handlers.AppRoutes())
}
