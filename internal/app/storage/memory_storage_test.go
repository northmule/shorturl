package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"testing"
	"time"
)

func TestMemoryStorage_StorageMethods(t *testing.T) {

	storage := NewMemoryStorage()

	tests := []struct {
		name     string
		testData models.URL
		want     models.URL
	}{
		{
			name: "#1_добавление, поиск",
			testData: models.URL{
				ShortURL: "1111",
				URL:      "https://google.com",
			},
			want: models.URL{
				ShortURL: "1111",
				URL:      "https://google.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := storage.Add(tt.testData)
			if err != nil {
				t.Errorf("Add() error = %#v", err)
			}
			url, _ := storage.FindByURL(tt.want.URL)
			if url.ShortURL != tt.want.ShortURL {
				t.Errorf("Add() ShortURL = %v, want %v", url.ShortURL, tt.want.ShortURL)
			}
			url, _ = storage.FindByShortURL(tt.want.ShortURL)
			if url.URL != tt.want.URL {
				t.Errorf("Add() ShortURL = %v, want %v", url.URL, tt.want.URL)
			}
		})
	}
}

// TestMemoryStorage_concurrentAdd Не должен упасть с fatal error: concurrent map writes
func TestMemoryStorage_concurrentAdd(t *testing.T) {
	storage := NewMemoryStorage()

	for i := 0; i < 200; i++ {
		go func() {
			storage.Add(models.URL{ShortURL: fmt.Sprintf("text%d", i), URL: "https://ya.ru"})
		}()
	}

	time.Sleep(time.Millisecond * 100)
	storage.Add(models.URL{ShortURL: "endKey", URL: "https://ya.ru"})
	if _, ok := (*storage.db)["endKey"]; !ok {
		t.Errorf("expected 'endKey' to be in the map")
	}
}
func TestMemoryStorage_CreateUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := models.User{
		ID:       1,
		Name:     "name",
		Login:    "Login",
		Password: "Login",
		UUID:     uuid.NewString(),
	}
	id, _ := storage.CreateUser(user)
	if id != int64(user.ID) {
		t.Errorf("CreateUser() id = %v, want %v", id, user.ID)
	}
}
