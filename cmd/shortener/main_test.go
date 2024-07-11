package main

import (
	"encoding/json"
	appStorage "github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"os"
	"testing"
)

type demoData []models.URL

func Test_restoreStorageData(t *testing.T) {
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
			t.Error(err)
			return
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

		storage := appStorage.NewStorage()
		restoreStorageData(file.Name(), storage)
		for _, url := range demoURLs {
			modelURL, err := storage.FindByURL(url.URL)
			if err != nil {
				t.Error(err)
			}
			if modelURL == nil {
				t.Error("Значений не найдено")
			}
		}
	})
}
