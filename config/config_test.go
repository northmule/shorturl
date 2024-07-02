package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "#1 задан флаг",
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
			name:  "#2 значения по умолчанию",
			env:   map[string]string{},
			flags: map[string]string{},
			want: Config{
				ServerURL:    ":8080",
				BaseShortURL: "http://localhost:8080",
			},
		},
		{
			name: "#3 есть env, нет флага",
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
			name: "#4 есть env, есть флаги",
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
				require.NoError(t, err)
			}

			for k, v := range tt.flags {
				args = append(args, fmt.Sprintf("-%s=%s", k, v))
			}
			os.Args = args

			configInit := Init()
			configInit.InitEnvConfig()
			configInit.InitFlagConfig()

			assert.Equal(t, tt.want.ServerURL, AppConfig.ServerURL)
			assert.Equal(t, tt.want.BaseShortURL, AppConfig.BaseShortURL)
		})
	}
}
