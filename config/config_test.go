package config

import (
	"fmt"
	"os"
	"testing"
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
			NewConfig()

			if tt.want.ServerURL != AppConfig.ServerURL {
				t.Errorf("Ожидается %#v пришло %#v", tt.want.ServerURL, AppConfig.ServerURL)
			}

			if tt.want.BaseShortURL != AppConfig.BaseShortURL {
				t.Errorf("Ожидается %#v пришло %#v", tt.want.BaseShortURL, AppConfig.BaseShortURL)
			}
		})
	}
}
