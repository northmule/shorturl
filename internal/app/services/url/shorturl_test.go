package url

import (
	"fmt"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"testing"
)

// storageMock структура хранилища
type storageMock struct {
	db *map[string]models.URL
}

// Add добавление нового значения
func (s *storageMock) Add(url models.URL) error {
	data := *s.db
	data[url.ShortURL] = url
	return nil
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
	return nil, fmt.Errorf("the url link was not found")
}

func TestShortURLService_DecodeURL(t *testing.T) {
	storageMock := &storageMock{
		db: &map[string]models.URL{},
	}
	NewShortURLService(storageMock)

	type fields struct {
		Storage      repositoryURLInterface
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
				Storage:      storageMock,
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
				Storage:      tt.fields.Storage,
				shortURLData: tt.fields.shortURLData,
			}
			shortURLResult, err := s.DecodeURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			modelURL, _ := s.Storage.FindByShortURL(shortURLResult.ShortURL)
			if modelURL.URL != tt.args.url {
				t.Errorf("DecodeURL() got = %v, want %v", modelURL.URL, tt.args.url)
			}

			modelURL, _ = s.Storage.FindByURL(tt.args.url)

			if modelURL.ShortURL != shortURLResult.ShortURL {
				t.Errorf("DecodeURL() got = %v, want %v", modelURL.ShortURL, shortURLResult.ShortURL)
			}
		})
	}
}

func TestShortURLService_EncodeShortURL(t *testing.T) {
	storageMock := &storageMock{
		db: &map[string]models.URL{},
	}
	NewShortURLService(storageMock)

	_ = storageMock.Add(models.URL{
		ShortURL: "123",
		URL:      "https://example.ru",
	})

	type fields struct {
		Storage      repositoryURLInterface
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
				Storage:      storageMock,
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
				Storage:      tt.fields.Storage,
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
