package main

import (
	"flag"
	"github.com/northmule/shorturl/config"
)

// parseFlags разбор аргументов коммандной строки
// Заполняет конфигурация приложения
func parseFlags() {
	flag.StringVar(&config.AppConfig.ServerURL, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.AppConfig.BaseShortURL, "b", "http://localhost:8000/", "the base address of the resulting shortened URL")
	flag.Parse()
}
