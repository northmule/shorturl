package config

import (
	"encoding/json"
	"errors"
	"os"
)

// JSONConfig Конфигурация приложения через JSON
type JSONConfig struct {
	path string
}

// JSONConfigFile Структура файла конфигурацииы
type JSONConfigFile struct {
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
}

// NewJSONConfig конструктор
func NewJSONConfig(path string) *JSONConfig {
	return &JSONConfig{
		path: path,
	}
}

// Init Инициализация данных JSON конфигурации
func (cfg *JSONConfig) Init(appConfig *Config) error {

	if cfg.path == "" {
		return errors.New("config file path is empty")
	}

	var JSONCfg JSONConfigFile
	var err error

	fb, err := os.ReadFile(cfg.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fb, &JSONCfg)
	if err != nil {
		return err
	}
	if appConfig.ServerURL == "" {
		appConfig.ServerURL = JSONCfg.ServerAddress
	}

	if appConfig.BaseShortURL == "" {
		appConfig.BaseShortURL = JSONCfg.BaseURL
	}

	if appConfig.FileStoragePath == "" {
		appConfig.FileStoragePath = JSONCfg.FileStoragePath
	}

	if appConfig.DataBaseDsn == "" {
		appConfig.DataBaseDsn = JSONCfg.DatabaseDSN
	}

	if !appConfig.EnableHTTPS {
		appConfig.EnableHTTPS = JSONCfg.EnableHTTPS
	}

	return nil
}
