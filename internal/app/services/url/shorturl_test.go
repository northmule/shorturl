package url

import (
	"fmt"
	"strings"
	"testing"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// storageMock структура хранилища
type storageMock struct {
	db *map[string]models.URL
}

// Add добавление нового значения
func (s *storageMock) Add(url models.URL) (int64, error) {
	data := *s.db
	data[url.ShortURL] = url
	return 0, nil
}

// FindByShortURL поиск по короткой ссылке
func (s *storageMock) FindByShortURL(shortURL string) (*models.URL, error) {
	data := *s.db
	if url, ok := data[shortURL]; ok {
		return &url, nil
	}

	return nil, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL
func (s *storageMock) FindByURL(url string) (*models.URL, error) {
	for _, modelURL := range *s.db {
		if modelURL.URL == url {
			return &modelURL, nil
		}
	}
	return new(models.URL), nil
}

func (s *storageMock) Ping() error {
	return nil
}

func (s *storageMock) MultiAdd(urls []models.URL) error {
	return nil
}

func (s *storageMock) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (s *storageMock) LikeURLToUser(urlID int64, userUUID string) error {
	return nil
}

func (s *storageMock) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (s *storageMock) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

func TestShortURLService_DecodeURL(t *testing.T) {
	_ = logger.InitLogger("fatal")
	storageMockInstance := &storageMock{
		db: &map[string]models.URL{},
	}
	NewShortURLService(storageMockInstance, storageMockInstance)

	type fields struct {
		Storage      storage.StorageQuery
		shortURLData ShortURLData
	}
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData *ShortURLData
		wantErr  bool
	}{
		{
			name: "#1_передать_url_получить_короткую_строку",
			fields: fields{
				Storage:      storageMockInstance,
				shortURLData: ShortURLData{},
			},
			args: args{
				url: "https://example.ru",
			},
			wantData: &ShortURLData{
				URL:      "https://example.ru",
				ShortURL: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{
				Finder:       tt.fields.Storage,
				Setter:       tt.fields.Storage,
				shortURLData: tt.fields.shortURLData,
			}
			shortURLResult, err := s.DecodeURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			modelURL, _ := s.Finder.FindByShortURL(shortURLResult.ShortURL)
			if modelURL.URL != tt.args.url {
				t.Errorf("DecodeURL() got = %v, want %v", modelURL.URL, tt.args.url)
			}

			modelURL, _ = s.Finder.FindByURL(tt.args.url)

			if modelURL.ShortURL != shortURLResult.ShortURL {
				t.Errorf("DecodeURL() got = %v, want %v", modelURL.ShortURL, shortURLResult.ShortURL)
			}
		})
	}
}

func TestShortURLService_EncodeShortURL(t *testing.T) {
	_ = logger.InitLogger("fatal")
	storageMockInstance := &storageMock{
		db: &map[string]models.URL{},
	}
	NewShortURLService(storageMockInstance, storageMockInstance)

	_, _ = storageMockInstance.Add(models.URL{
		ShortURL: "123",
		URL:      "https://example.ru",
	})

	type fields struct {
		Storage      storage.StorageQuery
		shortURLData ShortURLData
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantData *ShortURLData
		wantErr  bool
	}{
		{
			name: "#1_передать_короткую_ссылку_получить_url",
			fields: fields{
				Storage:      storageMockInstance,
				shortURLData: ShortURLData{},
			},
			args: args{
				shortURL: "123",
			},
			wantData: &ShortURLData{
				URL: "https://example.ru",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{
				Finder:       tt.fields.Storage,
				Setter:       tt.fields.Storage,
				shortURLData: tt.fields.shortURLData,
			}
			shortURLResult, err := s.EncodeShortURL(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeShortURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantData.URL != shortURLResult.URL {
				t.Errorf("EncodeShortURL() got = %v, want %v", shortURLResult.URL, tt.wantData.URL)
			}
		})
	}
}

func TestShortURLService_DecodeURLs(t *testing.T) {
	_ = logger.InitLogger("fatal")
	storageMockInstance := storage.NewMemoryStorage()

	tests := []struct {
		name    string
		Storage storage.StorageQuery
		urls    []string
		wantErr bool
	}{
		{
			name:    "#1_много_url",
			Storage: storageMockInstance,
			urls: []string{
				"https://habr.com/ru/feed/",
				"https://habr.com/ru/companies/gazprombank/articles/832810/",
				"https://habr.com/ru/companies/f_a_c_c_t/news/833140/",
				"https://habr.com/ru/news/833130/",
				"https://habr.com/ru/companies/kts/news/833080/",
				"https://habr.com/ru/companies/aenix/news/833030/",
				"https://habr.com/ru/companies/alfa/articles/",
				"https://habr.com/ru/companies/otus/articles/",
			},
			wantErr: false,
		},
		{
			name:    "#2_пустой_список",
			Storage: storageMockInstance,
			urls:    []string{},
			wantErr: false,
		},
		{
			name:    "#3_дубли",
			Storage: storageMockInstance,
			urls: []string{
				"https://habr.com/ru/feed/",
				"https://habr.com/ru/feed/",
				"https://habr.com/ru/feed/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShortURLService{
				Finder:       tt.Storage,
				Setter:       tt.Storage,
				shortURLData: ShortURLData{},
			}
			_, err := s.DecodeURLs(tt.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, url := range tt.urls {
				_, err := s.Finder.FindByURL(url)
				if err != nil {
					t.Errorf("DecodeURL() error = %v", err)
				}
			}
		})
	}
}

func BenchmarkNewRandomString(b *testing.B) {
	_ = logger.InitLogger("fatal")
	for i := 0; i < b.N; i++ {
		newRandomString(ShortURLDefaultSize)
	}
}

func BenchmarkDecodeURLs(b *testing.B) {
	_ = logger.InitLogger("fatal")
	storageMemoryMock := storage.NewMemoryStorage()
	service := &ShortURLService{
		Finder:       storageMemoryMock,
		Setter:       storageMemoryMock,
		shortURLData: ShortURLData{},
	}
	testData := strings.Repeat("A ", 100)
	urls := strings.Split(testData, " ")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = service.DecodeURLs(urls)
	}
}
