package storage

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_StorageMethods(t *testing.T) {
	_ = logger.InitLogger("fatal")
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
	_ = logger.InitLogger("fatal")
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

func TestMemoryStorage_FindUserByLoginAndPasswordHash(t *testing.T) {
	_ = logger.InitLogger("fatal")
	storage := NewMemoryStorage()
	newUser := models.User{
		Name:     "name",
		Login:    "Login",
		Password: "Password",
		UUID:     uuid.NewString(),
	}
	_, _ = storage.CreateUser(newUser)
	user, _ := storage.FindUserByLoginAndPasswordHash("Login", "Password")
	if user.Login != "Login" {
		t.Errorf("FindUserByLoginAndPasswordHash() Login = %v, want %v", user.Login, "Login")
	}
}

func TestMemoryStorage__FindUrlsByUserID(t *testing.T) {
	_ = logger.InitLogger("fatal")
	storage := NewMemoryStorage()
	userUUID := "1111-2222-33333-44444"
	_, _ = storage.CreateUser(models.User{
		Name:     "name",
		Login:    "Login",
		Password: "Password",
		UUID:     userUUID,
	})
	urlID, _ := storage.Add(models.URL{
		ShortURL: "qqwww",
		URL:      "https://google.com",
	})
	_ = storage.LikeURLToUser(urlID, userUUID)

	userURLs, _ := storage.FindUrlsByUserID(userUUID)
	if len(*userURLs) == 0 {
		t.Errorf("FindUrlsByUserID() userURLs = %v, want %v", 0, 1)
	}
	var isExist bool
	for _, u := range *userURLs {
		if u.ShortURL == "qqwww" {
			isExist = true
		}
	}
	if !isExist {
		t.Errorf("FindUrlsByUserID() userURLs = %v, want %v", 0, 1)
	}
}

func BenchmarkMemoryStorage(b *testing.B) {
	_ = logger.InitLogger("fatal")

	b.Run("поиск_ссылок_созданных_пользователем", func(b *testing.B) {
		storage := NewMemoryStorage()
		userUUID := "1111-2222-33333-44444"
		_, _ = storage.CreateUser(models.User{
			Name:     "name",
			Login:    "Login",
			Password: "Password",
			UUID:     userUUID,
		})
		urlID, _ := storage.Add(models.URL{
			ShortURL: "qqwww",
			URL:      "https://google.com",
		})
		_ = storage.LikeURLToUser(urlID, userUUID)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userURLs, _ := storage.FindUrlsByUserID(userUUID)
			if len(*userURLs) == 0 {
				b.Errorf("URL-s не найдены")
			}
		}
	})

	b.Run("поиск_пользователя_по_логину_и_паролю", func(b *testing.B) {
		storage := NewMemoryStorage()
		userUUID := "1111-2222-33333-44444"
		_, _ = storage.CreateUser(models.User{
			Name:     "name",
			Login:    "Login",
			Password: "Password",
			UUID:     userUUID,
		})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			user, _ := storage.FindUserByLoginAndPasswordHash("Login", "Password")
			if user.Login != "Login" {
				b.Errorf("Логин не совпадает")
			}
		}
	})

	b.Run("создание_нового_пользователя", func(b *testing.B) {
		storage := NewMemoryStorage()
		b.ResetTimer()
		testID := int64(0)
		for i := 0; i < b.N; i++ {
			userID, _ := storage.CreateUser(models.User{
				Name:     "name",
				Login:    "Login",
				Password: "Password",
				UUID:     "1111-2222-33333-44444",
			})
			testID++
			if userID != testID {
				b.Errorf("ID не совпадаеют")
			}
		}
	})

	b.Run("поиск_по_url", func(b *testing.B) {
		storage := NewMemoryStorage()
		storage.Add(models.URL{
			ShortURL: "111fghfhfgh1",
			URL:      "https://google.com",
		})

		for i := 1; i < 100000; i++ {
			storage.Add(models.URL{
				ShortURL: "asdfsfadf",
				URL:      "https://habr.ru/news_" + string(rune(i)),
			})
		}

		storage.Add(models.URL{
			ShortURL: "2222vbxcbcvbc2",
			URL:      "https://ya.ru",
		})
		var url *models.URL
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			url, _ = storage.FindByURL("https://google.com")
			if url == nil {
				b.Errorf("URL не найден")
			}
			url, _ = storage.FindByURL("https://ya.ru")
			if url == nil {
				b.Errorf("URL не найден")
			}
		}
	})

	b.Run("вставка_нескольких_значений", func(b *testing.B) {
		storage := NewMemoryStorage()

		urls := make([]models.URL, 1000)
		for i := 0; i < 1000; i++ {
			urls[i] = models.URL{
				ShortURL: "asdfsfadf",
				URL:      "https://habr.ru/news_" + strconv.Itoa(i),
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			storage.MultiAdd(urls)
		}
	})
}

func TestMemoryStorage_GetCountShortURL(t *testing.T) {
	storage := NewMemoryStorage()
	cnt, _ := storage.GetCountShortURL()
	assert.Equal(t, int64(1), cnt)
	storage.Add(models.URL{ShortURL: "123", URL: "https://ya.ru"})
	storage.Add(models.URL{ShortURL: "321", URL: "https://ya1.ru"})
	cnt, _ = storage.GetCountShortURL()
	assert.Equal(t, int64(3), cnt)
}

func TestMemoryStorage_GetCountUser(t *testing.T) {
	storage := NewMemoryStorage()
	user := models.User{
		ID:       1,
		Name:     "name",
		Login:    "Login",
		Password: "Login",
		UUID:     uuid.NewString(),
	}
	_, _ = storage.CreateUser(user)
	cnt, _ := storage.GetCountUser()
	assert.Equal(t, int64(1), cnt)
}
