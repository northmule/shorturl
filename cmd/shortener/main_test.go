package main

import (
	"encoding/json"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"os"
	"testing"
)

type demoData []models.URL

func Test_restoreStorageData(t *testing.T) {
	_ = logger.NewLogger("fatal")

	t.Run("првоерка_восстановления_значений_из_файла", func(t *testing.T) {
		demoURLs := demoData{
			{
				ID:       1,
				URL:      "111",
				ShortURL: "1",
			},
			{
				ID:       2,
				URL:      "222",
				ShortURL: "2",
			},
			{
				ID:       1,
				URL:      "333",
				ShortURL: "3",
			},
		}

		file, err := os.CreateTemp("/tmp", "gotest_*_.json")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(file.Name())
		jsonEncoder := json.NewEncoder(file) // будет записанно в файл
		for _, v := range demoURLs {
			err = jsonEncoder.Encode(v)
			if err != nil {
				t.Error(err)
				continue
			}
		}
		config.AppConfig.FileStoragePath = file.Name()
		storage := appStorage.NewStorage(true)

		for _, url := range demoURLs {
			modelURL, err := storage.FindByURL(url.URL)
			if err != nil {
				t.Error(err)
			}
			if modelURL == nil {
				t.Errorf("Значений не найдено: storage.FindByURL(%s)", url.URL)
			}

			modelURL, err = storage.FindByShortURL(url.ShortURL)
			if err != nil {
				t.Error(err)
			}
			if modelURL == nil {
				t.Errorf("Значений не найдено: storage.FindByShortURL(%s)", url.ShortURL)
			}
		}
	})

	t.Run("путь_к_файлу_не_передан", func(t *testing.T) {
		demoURLs := demoData{
			{
				ID:       1,
				URL:      "111",
				ShortURL: "1",
			},
			{
				ID:       2,
				URL:      "222",
				ShortURL: "2",
			},
			{
				ID:       1,
				URL:      "333",
				ShortURL: "3",
			},
		}
		file, err := os.CreateTemp("/tmp", "gotest_*_.json")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(file.Name())
		jsonEncoder := json.NewEncoder(file)
		for _, v := range demoURLs {
			err = jsonEncoder.Encode(v)
			if err != nil {
				t.Error(err)
				continue
			}
		}

		storage := appStorage.NewStorage(true)
		for _, url := range demoURLs {
			modelURL, _ := storage.FindByURL(url.URL)
			if modelURL != nil {
				t.Errorf("Значений не ожидалось: storage.FindByURL(%s)", url.URL)
			}
		}
	})
}
