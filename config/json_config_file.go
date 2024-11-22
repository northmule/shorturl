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

	var JSONCfg ConfigurationFile
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

	if appConfig.TrustedSubnet == "" {
		appConfig.TrustedSubnet = JSONCfg.TrustedSubnet
	}

	return nil
}
