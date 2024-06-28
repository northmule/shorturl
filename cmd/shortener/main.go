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
func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.DecodeHandler)
	mux.HandleFunc(`/{id}`, handlers.EncodeHandler)

	return http.ListenAndServe(configs.ServerURL, mux)
}
