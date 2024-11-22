package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name          string
		want          Config
		actual        string
		expectedError bool
	}{
		{
			name: "валидный_конфиг",
			want: Config{
				ServerURL:       "localhost:8080",
				BaseShortURL:    "http://localhost:8080",
				FileStoragePath: "/tmp/storage",
				DataBaseDsn:     "/dbname",
				EnableHTTPS:     true,
				TrustedSubnet:   "192.168.0.1/24",
			},
			actual: `{
		"server_address": "localhost:8080",
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true,
		"trusted_subnet": "192.168.0.1/24"
	}`,
		},
		{
			name: "не_валидный_конфиг",
			want: Config{},
			actual: `{
		"server_address": "localhost:8080":,
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true
	}`,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "config.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())
			_, err = tmpFile.Write([]byte(tt.actual))
			if err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			err = tmpFile.Close()
			if err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			cfg := NewJSONConfig(tmpFile.Name())
			actualConfig := Config{}

			err = cfg.Init(&actualConfig)
			if tt.expectedError {
				if err == nil {
					t.Errorf("Init() expected error but got nil")
				}
			}
			if !tt.expectedError {
				if err != nil {
					t.Errorf("Init() expected no error but got %v", err)
				}
			}
			if diff := cmp.Diff(tt.want, actualConfig); diff != "" {
				t.Errorf("Config mismatch (-expected +got):\n%s", diff)
			}

		})
	}
}

func TestInit_Path(t *testing.T) {
	tests := []struct {
		name          string
		path          func() string
		expectedError bool
	}{
		{
			name: "не_указан_путь_к_файлу",
			path: func() string {
				return ""
			},
			expectedError: true,
		},
		{
			name: "путь_не_существует",
			path: func() string {
				return "karamba.json"
			},
			expectedError: true,
		},
		{
			name: "фай_указан_и_существует",
			path: func() string {
				fileData := `{
		"server_address": "localhost:8080",
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true
	}`
				pathFile, _ := os.CreateTemp("", "config.json")
				_, _ = pathFile.Write([]byte(fileData))
				return pathFile.Name()
			},
			expectedError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathFile := tt.path()
			cfg := NewJSONConfig(pathFile)
			appConfig := &Config{}
			if pathFile != "" {
				defer os.Remove(pathFile)
			}
			err := cfg.Init(appConfig)
			if tt.expectedError && err == nil {
				t.Errorf("Init() expected error but got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Init() expected no error but got %v", err)
			}
		})
	}
}
