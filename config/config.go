package config

// Config Конфигурация приложения
type Config struct {
	// Адрес сервера и порт
	ServerURL string
	// Базовый адрес результирующего сокращённого URL
	//(значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
	BaseShortURL string
}

// AppConfig глобальная переменная конфигурации
var AppConfig Config

// Init Инициализация конфигурации приложения
func Init() {
	AppConfig = Config{}
}
