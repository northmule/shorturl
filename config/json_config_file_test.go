package config

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	configContent := `{
		"server_address": "localhost:8080",
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg := NewJSONConfig(tmpFile.Name())
	appConfig := &Config{}

	err = cfg.Init(appConfig)
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if appConfig.ServerURL != "localhost:8080" {
		t.Errorf("Expected ServerURL to be 'localhost:8080', but got '%s'", appConfig.ServerURL)
	}
	if appConfig.BaseShortURL != "http://localhost:8080" {
		t.Errorf("Expected BaseShortURL to be 'http://localhost:8080', but got '%s'", appConfig.BaseShortURL)
	}
	if appConfig.FileStoragePath != "/tmp/storage" {
		t.Errorf("Expected FileStoragePath to be '/tmp/storage', but got '%s'", appConfig.FileStoragePath)
	}
	if appConfig.DataBaseDsn != "/dbname" {
		t.Errorf("Expected DataBaseDsn to be '/dbname', but got '%s'", appConfig.DataBaseDsn)
	}
	if !appConfig.EnableHTTPS {
		t.Errorf("Expected EnableHTTPS to be true, but got false")
	}
}

func TestInit_EmptyPath(t *testing.T) {
	cfg := NewJSONConfig("")
	appConfig := &Config{}

	err := cfg.Init(appConfig)
	if err == nil {
		t.Error("Expected error for empty config file path, but got nil")
	}

}

func TestInit_InvalidJSON(t *testing.T) {
	configContent := `{
		"server_address": "localhost:8080"///,
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "user:pass@tcp(localhost:3306)/dbname",
		"enable_https": true,
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte(configContent))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	cfg := NewJSONConfig(tmpFile.Name())
	appConfig := &Config{}

	err = cfg.Init(appConfig)
	if err == nil {
		t.Error("Expected error for invalid JSON, but got nil")
	}
}
func TestInit_FileNotFound(t *testing.T) {
	cfg := NewJSONConfig("nonexistent.json")
	appConfig := &Config{}

	err := cfg.Init(appConfig)
	if err == nil {
		t.Error("Expected error for file not found, but got nil")
	}
}
