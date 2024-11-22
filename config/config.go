package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/caarlos0/env"
)

// Приоритет параметров сервера должен быть таким:
// Если указана переменная окружения, то используется она.
// Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
// Если нет ни переменной окружения, ни флага, ни json конфигурации -  то используется значение по умолчанию.

// Параметры по умолчанию.
const (
	addressAndPortDefault     = ":8080"
	baseAddressDefault        = "http://localhost:8080"
	pathFileStorage           = "/tmp/short-url-db.json"
	DataBaseConnectionTimeOut = 10000
	pprofEnabledDefault       = true
	enableHTTPSDefault        = false
)

// Config Конфигурация приложения.
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
	// Аткивация pprof
	PprofEnabled bool `env:"PPROF_ENABLED"`
	// Запуск сервера https
	EnableHTTPS bool `env:"ENABLE_HTTPS"`
	// Путь к файлу конфигурации приложения
	Config string `env:"CONFIG"`
	// Доверенная сеть
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

// ConfigurationFile Структура файла конфигурацииы
type ConfigurationFile struct {
	// ServerAddress аналог переменной окружения SERVER_ADDRESS или флага -a
	ServerAddress string `json:"server_address"`
	// BaseURL аналог переменной окружения BASE_URL или флага -b
	BaseURL string `json:"base_url"`
	// FileStoragePath аналог переменной окружения FILE_STORAGE_PATH или флага -f
	FileStoragePath string `json:"file_storage_path"`
	// DatabaseDSN аналог переменной окружения DATABASE_DSN или флага -d
	DatabaseDSN string `json:"database_dsn"`
	// EnableHTTPS аналог переменной окружения ENABLE_HTTPS или флага -s
	EnableHTTPS bool `json:"enable_https"`
	// Доверенная сеть
	TrustedSubnet string `json:"trusted_subnet"`
}

// InitConfig инициализация настроек приложения.
type InitConfig interface {
	InitEnvConfig() error
	InitFlagConfig() error
}

// AppConfig глобальная переменная конфигурации.
var AppConfig Config

// NewConfig Инициализация конфигурации приложения.
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
	err = AppConfig.InitJSONConfig()
	if err != nil {
		return nil, err
	}

	AppConfig.initDefaultConfig()
	return &AppConfig, nil
}

// InitEnvConfig разбор настроек из env.
func (c *Config) InitEnvConfig() error {
	return initEnvConfig(c)
}

// InitFlagConfig разбор настроек из флагов.
func (c *Config) InitFlagConfig() error {
	return initFlagConfig(c)
}

// initEnvConfig прасинг env переменных.
func initEnvConfig(appConfig *Config) error {
	// Заполнение значений из окружения
	err := env.Parse(appConfig)
	if err != nil {
		return err
	}
	return nil
}

// InitJSONConfig данные из JSON конфигурации (самый низкий приоритет)
func (c *Config) InitJSONConfig() error {
	if c.Config == "" {
		return nil
	}
	cfg := NewJSONConfig(c.Config)
	return cfg.Init(c)
}

// initFlagConfig парсинг флагов командной строки.
func initFlagConfig(appConfig *Config) error {
	// На каждый новый запуск новая структура флагов
	configFlag := flag.FlagSet{}
	flagServerURLValue := configFlag.String("a", "", "address and port to run server")
	flagBaseShortURLValue := configFlag.String("b", "", "the base address of the resulting shortened URL")
	// Если указан пустой флга, запись в файл отключается
	flagFileStoragePath := configFlag.String("f", "", "the path to the file for storing links")
	// Строка подключения базы данных
	flagDataBaseDsn := configFlag.String("d", "", "specify the connection string to the database")
	pprofEnabled := configFlag.Bool("pprof", false, "enable pprof")
	enableHTTPS := configFlag.Bool("s", false, "launching the https server")

	flagFileConfigShortApp := configFlag.String("c", "", "the path to the application configuration file")
	flagFileConfigFullApp := configFlag.String("config", "", "the path to the application configuration file")
	flagFileConfigTrustedSubnet := configFlag.String("t", "", "trusted network")

	err := configFlag.Parse(os.Args[1:])
	if err != nil {
		return errors.Join(errors.New("failed to parse flags"), err)
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
	appConfig.PprofEnabled = *pprofEnabled
	appConfig.EnableHTTPS = *enableHTTPS

	var flagFileConfigApp string
	if *flagFileConfigShortApp != "" {
		flagFileConfigApp = *flagFileConfigShortApp
	}
	if *flagFileConfigFullApp != "" {
		flagFileConfigApp = *flagFileConfigFullApp
	}

	if appConfig.Config == "" {
		appConfig.Config = flagFileConfigApp
	}

	if *flagFileConfigTrustedSubnet == "" {
		appConfig.TrustedSubnet = *flagFileConfigTrustedSubnet
	}

	appConfig.DataBaseDsn = strings.ReplaceAll(appConfig.DataBaseDsn, "\"", "")
	return nil
}

func (c *Config) initDefaultConfig() {
	if c.ServerURL == "" {
		c.ServerURL = addressAndPortDefault
	}

	if c.BaseShortURL == "" {
		c.BaseShortURL = baseAddressDefault
	}

	if c.FileStoragePath == "" {
		c.FileStoragePath = pathFileStorage
	}

	if !c.PprofEnabled {
		c.PprofEnabled = pprofEnabledDefault
	}

	if !c.EnableHTTPS {
		c.EnableHTTPS = enableHTTPSDefault
	}
}
