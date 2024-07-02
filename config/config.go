package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"os"
)

// Приоритет параметров сервера должен быть таким:
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

const addressAndPortDefault = ":8080"
const baseAddressDefault = "http://localhost:8080"

// Config Конфигурация приложения
type Config struct {
	// Адрес сервера и порт
	ServerURL string `env:"SERVER_ADDRESS"`
	// Базовый адрес результирующего сокращённого URL
	//(значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
	BaseShortURL string `env:"BASE_URL"`
	DatabasePath string
}

type InitConfig interface {
	InitEnvConfig()
	InitFlagConfig()
	InitStaticConfig()
}

// AppConfig глобальная переменная конфигурации
var AppConfig Config

// Init Инициализация конфигурации приложения
func Init() InitConfig {
	AppConfig = Config{}
	return &AppConfig

}

func (c *Config) InitEnvConfig() {
	initEnvConfig(c)
}

func (c *Config) InitFlagConfig() {
	initFlagConfig(c)
}

func (c *Config) InitStaticConfig() {
	initStaticConfig(c)
}

// initEnvConfig прасинг env переменных
func initEnvConfig(appConfig *Config) {
	// Заполнение значений из окружения
	err := env.Parse(appConfig)
	if err != nil {
		_ = fmt.Errorf("parse error env: %s", err)
	}
}

// initFlagConfig парсинг флагов командной строки
func initFlagConfig(appConfig *Config) {
	// На каждый новый запуск новая структура флагов
	configFlag := flag.FlagSet{}
	flagServerURLValue := configFlag.String("a", addressAndPortDefault, "address and port to run server")
	flagBaseShortURLValue := configFlag.String("b", baseAddressDefault, "the base address of the resulting shortened URL")
	err := configFlag.Parse(os.Args[1:])
	if err != nil {
		_ = fmt.Errorf("parse error os.Args: %s", err)
	}

	if appConfig.ServerURL == "" {
		appConfig.ServerURL = *flagServerURLValue
	}
	if appConfig.BaseShortURL == "" {
		appConfig.BaseShortURL = *flagBaseShortURLValue
	}

}

// initStaticConfig todo убрать в будущем
func initStaticConfig(appConfig *Config) {
	appConfig.DatabasePath = "/home/djo/GolandProjects/shorturl/shorturl.db"
}
