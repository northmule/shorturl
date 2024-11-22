package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/assert"
)

type demoData []models.URL

func TestFileStorage_restoreStorageData(t *testing.T) {
	_ = logger.InitLogger("fatal")

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
		file.Close()

		fileStorage, err := os.OpenFile(file.Name(), os.O_RDWR, 0666)
		if err != nil {
			t.Fatal(err)
		}
		storage := NewFileStorage(fileStorage)

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

}

func TestFileStorage_Add(t *testing.T) {
	_ = logger.InitLogger("fatal")

	t.Run("Добавление_значения", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "TestFileStorage_Add_*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())
		fileStorage := NewFileStorage(file)

		url := models.URL{
			ID:       1,
			ShortURL: "aaa",
			URL:      "bbbbbbb",
		}
		_, err = fileStorage.Add(url)
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}
	})

	t.Run("Проверка_добавленного_значения", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "TestFileStorage_Add_*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())
		fileStorage := NewFileStorage(file)

		url := models.URL{
			ID:       1,
			ShortURL: "aaa",
			URL:      "bbbbbbb",
		}
		_, _ = fileStorage.Add(url)
		findValue, err := fileStorage.FindByURL(url.URL)
		if findValue == nil {
			t.Errorf("FindByURL() error = %v", err)
		}
	})

	t.Run("Запись_множества_значений_с_проверкой_наличия_одного_значения", func(t *testing.T) {
		_ = logger.InitLogger("fatal")
		file, err := os.CreateTemp("/tmp", "TestFileStorage_Add_*.json")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())
		fileStorage := NewFileStorage(file)

		for i := 0; i < 200; i++ {
			go func() {
				fileStorage.Add(models.URL{
					ID:       uint(i),
					ShortURL: fmt.Sprintf("text%d", i),
					URL:      fmt.Sprintf("https://ya.ru/%d", i),
				})
			}()
		}

		time.Sleep(time.Millisecond * 100)
		_, err = fileStorage.Add(models.URL{ShortURL: "endKey", URL: "https://ya.ru"})
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}
		findValue, err := fileStorage.FindByURL("https://ya.ru")
		if findValue == nil {
			t.Errorf("FindByURL() error = %v", err)
		}
	})

}

func TestCreateUser(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test-storage-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	storage := NewFileStorage(tempFile)
	if storage == nil {
		t.Fatalf("Failed to initialize FileStorage")
	}
	defer storage.Close()

	user := models.User{
		Login: "testuser",
		UUID:  "testuuid",
	}

	_, err = storage.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser returned an error: %v", err)
	}

	tempFile.Close()
	tempFile, err = os.Open(tempFile.Name() + "user.json")
	if err != nil {
		t.Fatalf("Failed to open user file: %v", err)
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(tempFile)
	if !scanner.Scan() {
		t.Fatalf("Expected at least one line in the user file")
	}
	line := scanner.Text()
	var storedUser models.User
	err = json.Unmarshal([]byte(line), &storedUser)
	if err != nil {
		t.Errorf("Failed to unmarshal user data: %v", err)
	}
	if storedUser.Login != user.Login || storedUser.UUID != user.UUID {
		t.Errorf("Stored user data does not match expected data")
	}
}

func TestGetCountShortURL(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test-storage-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	storage := NewFileStorage(tempFile)
	if storage == nil {
		t.Fatalf("Failed to initialize FileStorage")
	}
	defer storage.Close()
	tempFile.Close()
	url := models.URL{
		ID:       1,
		ShortURL: "aaa",
		URL:      "bbbbbbb",
	}
	_, _ = storage.Add(url)

	cnt, _ := storage.GetCountShortURL()

	assert.Equal(t, int64(1), cnt)
}

func TestGetCountUser(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test-storage-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	storage := NewFileStorage(tempFile)
	if storage == nil {
		t.Fatalf("Failed to initialize FileStorage")
	}
	defer storage.Close()

	user := models.User{
		Login: "testuser",
		UUID:  "testuuid",
	}

	_, err = storage.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser returned an error: %v", err)
	}

	cnt, _ := storage.GetCountUser()
	assert.Equal(t, int64(1), cnt)
}
