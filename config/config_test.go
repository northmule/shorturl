package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAFlagAndBFlag(t *testing.T) {
	tests := []struct {
		name  string
		env   map[string]string
		flags map[string]string
		want  Config
	}{
		{
			name: "#1_задан_флаг",
			env:  map[string]string{},
			flags: map[string]string{
				"a": ":8082",
				"b": "http://localhost:8000/",
			},
			want: Config{
				ServerURL:    ":8082",
				BaseShortURL: "http://localhost:8000/",
			},
		},
		{
			name:  "#2_значения_по_умолчанию",
			env:   map[string]string{},
			flags: map[string]string{},
			want: Config{
				ServerURL:    ":8080",
				BaseShortURL: "http://localhost:8080",
			},
		},
		{
			name: "#3_есть_env_нет_флага",
			env: map[string]string{
				"SERVER_ADDRESS": ":9082",
				"BASE_URL":       "http://localhost:8082/",
			},
			flags: map[string]string{},
			want: Config{
				ServerURL:    ":9082",
				BaseShortURL: "http://localhost:8082/",
			},
		},
		{
			name: "#4_есть_env_есть_флаги",
			env: map[string]string{
				"SERVER_ADDRESS": ":9082",
				"BASE_URL":       "http://localhost:8082/",
			},
			flags: map[string]string{
				"a": ":6544",
				"b": "http://localhost:8000/",
			},
			want: Config{
				ServerURL:    ":9082",
				BaseShortURL: "http://localhost:8082/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			var args []string
			args = append(args, "config_test.go")

			for k, v := range tt.env {
				err := os.Setenv(k, v)
				if err != nil {
					t.Error(err)
				}
			}

			for k, v := range tt.flags {
				args = append(args, fmt.Sprintf("-%s=%s", k, v))
			}
			os.Args = args

			// Инициализация после подготовки флагов и переменных
			_, err := NewConfig()
			if err != nil {
				t.Error(err)
			}

			if tt.want.ServerURL != AppConfig.ServerURL {
				t.Errorf("Ожидается %#v пришло %#v", tt.want.ServerURL, AppConfig.ServerURL)
			}

			if tt.want.BaseShortURL != AppConfig.BaseShortURL {
				t.Errorf("Ожидается %#v пришло %#v", tt.want.BaseShortURL, AppConfig.BaseShortURL)
			}
		})
	}
}

func TestInitJSONConfig(t *testing.T) {
	jsonConfig := `{
		"server_address": "localhost:8080",
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true
	}`
	jsonFile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(jsonFile.Name())

	_, err = jsonFile.WriteString(jsonConfig)
	assert.NoError(t, err)
	err = jsonFile.Close()
	assert.NoError(t, err)

	var config Config

	config.Config = jsonFile.Name()
	err = config.InitJSONConfig()
	if err != nil {
		return
	}
	assert.Equal(t, "localhost:8080", config.ServerURL)
	assert.Equal(t, "http://localhost:8080", config.BaseShortURL)
	assert.Equal(t, "/tmp/storage", config.FileStoragePath)
	assert.Equal(t, "/dbname", config.DataBaseDsn)
	assert.True(t, config.EnableHTTPS)
	assert.Equal(t, jsonFile.Name(), config.Config)
}

func TestNewConfig(t *testing.T) {

	_ = os.Setenv("SERVER_ADDRESS", "mocked_address")
	_ = os.Setenv("BASE_URL", "mocked_base_url")
	_ = os.Setenv("FILE_STORAGE_PATH", "mocked_file_path")
	_ = os.Setenv("DATABASE_DSN", "mocked_db_dsn")
	_ = os.Setenv("PPROF_ENABLED", "true")
	_ = os.Setenv("ENABLE_HTTPS", "true")

	jsonFile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)

	_ = os.Setenv("CONFIG", jsonFile.Name())
	defer os.Remove(jsonFile.Name())

	os.Args = []string{"cmd", "-a", "mocked_address_cmd", "-b", "mocked_base_url_cmd", "-f", "mocked_file_path_cmd", "-d", "mocked_db_dsn_cmd", "-pprof", "-s"}

	jsonConfig := `{
		"server_address": "localhost:8080",
		"base_url": "http://localhost:8080",
		"file_storage_path": "/tmp/storage",
		"database_dsn": "/dbname",
		"enable_https": true
	}`

	_, err = jsonFile.WriteString(jsonConfig)
	assert.NoError(t, err)
	err = jsonFile.Close()
	assert.NoError(t, err)

	config, err := NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, "mocked_address", config.ServerURL)
	assert.Equal(t, "mocked_base_url", config.BaseShortURL)
	assert.Equal(t, "mocked_file_path", config.FileStoragePath)
	assert.Equal(t, "mocked_db_dsn", config.DataBaseDsn)
	assert.True(t, config.PprofEnabled)
	assert.True(t, config.EnableHTTPS)
	assert.Equal(t, jsonFile.Name(), config.Config)
}
