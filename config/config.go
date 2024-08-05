package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/northmule/shorturl/internal/app/logger"
	"os"
	"strings"
)

// Приоритет параметров сервера должен быть таким:
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

const addressAndPortDefault = ":8080"
const baseAddressDefault = "http://localhost:8080"
const pathFileStorage = "/tmp/short-url-db.json"
const DataBaseConnectionTimeOut = 10

// Config Конфигурация приложения
type Config struct {
	// Адрес сервера и порт
	ServerURL string `env:"SERVER_ADDRESS"`
	// Базовый адрес результирующего сокращённого URL
	//(значение: адрес сервера перед коротким URL, например http://localhost:8000/qsd54gFg).
	BaseShortURL string `env:"BASE_URL"`
	// Путь для хранения ссылок
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// Строка подключения к БД
	DataBaseDsn string `env:"DATABASE_DSN"`
}

type InitConfig interface {
	InitEnvConfig() error
	InitFlagConfig() error
}

// AppConfig глобальная переменная конфигурации
var AppConfig Config

// NewConfig Инициализация конфигурации приложения
func NewConfig() (*Config, error) {
	AppConfig = Config{}
	err := AppConfig.InitEnvConfig()
	if err != nil {
		return nil, err
	}
	err = AppConfig.InitFlagConfig()
	if err != nil {
		return nil, err
	}
	return &AppConfig, nil
}

func (c *Config) InitEnvConfig() error {
	return initEnvConfig(c)
}

func (c *Config) InitFlagConfig() error {
	return initFlagConfig(c)
}

// initEnvConfig прасинг env переменных
func initEnvConfig(appConfig *Config) error {
	// Заполнение значений из окружения
	err := env.Parse(appConfig)
	if err != nil {
		return err
	}
	return nil
}

// initFlagConfig парсинг флагов командной строки
func initFlagConfig(appConfig *Config) error {
	// На каждый новый запуск новая структура флагов
	configFlag := flag.FlagSet{}
	flagServerURLValue := configFlag.String("a", addressAndPortDefault, "address and port to run server")
	flagBaseShortURLValue := configFlag.String("b", baseAddressDefault, "the base address of the resulting shortened URL")
	// Если указан пустой флга, запись в файл отключается
	flagFileStoragePath := configFlag.String("f", pathFileStorage, "the path to the file for storing links")
	// Строка подключения базы данных
	flagDataBaseDsn := configFlag.String("d", "", "specify the connection string to the database")

	err := configFlag.Parse(os.Args[1:])
	if err != nil {
		logger.LogSugar.Error("configFlag.Parse error", err)
		return err
	}

	if appConfig.ServerURL == "" {
		appConfig.ServerURL = *flagServerURLValue
	}
	if appConfig.BaseShortURL == "" {
		appConfig.BaseShortURL = *flagBaseShortURLValue
	}
	if appConfig.FileStoragePath == "" {
		appConfig.FileStoragePath = *flagFileStoragePath
	}
	if appConfig.DataBaseDsn == "" {
		appConfig.DataBaseDsn = *flagDataBaseDsn

	}
	appConfig.DataBaseDsn = strings.ReplaceAll(appConfig.DataBaseDsn, "\"", "")
	return nil
}
