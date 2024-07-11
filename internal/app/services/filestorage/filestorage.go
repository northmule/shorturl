package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"os"
)

type Setter struct {
	file            *os.File
	shortURLService *url.ShortURLService
	encoder         *json.Encoder
}

// NewSetter конструктор
func NewSetter(filename string, shortURLService *url.ShortURLService) (*Setter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Sugar().Errorw("Failed to open file", "filename", filename, "error", err)
		return nil, err
	}
	return &Setter{
		file:            file,
		shortURLService: shortURLService,
		encoder:         json.NewEncoder(file),
	}, nil
}

// WriteURL Запись стуктуры в файл
func (s *Setter) WriteURL(url models.URL) error {
	return s.encoder.Encode(url)
}

// DecodeURL добавляем запись для декодера URL
func (s *Setter) DecodeURL(url string) (*url.ShortURLData, error) {
	data, err := s.shortURLService.DecodeURL(url)
	if err != nil {
		return nil, err
	}
	modelURL := models.URL{
		URL:      data.URL,
		ShortURL: data.ShortURL,
	}
	err = s.WriteURL(modelURL)
	if err != nil {
		return nil, err
	}
	return data, err
}

// EncodeShortURL вызываем оригинальный энкодер
func (s *Setter) EncodeShortURL(url string) (*url.ShortURLData, error) {
	return s.shortURLService.EncodeShortURL(url)
}

// WriteURLs Запись мапы в файл
func (s *Setter) WriteURLs(data *map[string]models.URL) {
	for k, v := range *data {
		if err := s.WriteURL(v); err != nil {
			logger.Log.Sugar().Error("Failed to encode url", "url", k, "error", err)
			continue
		}
	}
}

func (s *Setter) Close() error {
	return s.file.Close()
}

type Getter struct {
	file    *os.File
	decoder *json.Decoder
	scanner *bufio.Scanner
}

func NewGetter(filename string) (*Getter, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Sugar().Errorw("Failed to open file", "filename", filename, "error", err)
		return nil, err
	}
	return &Getter{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// ReadURL чтение структуры из файла
func (g *Getter) ReadURL() (*models.URL, error) {
	url := &models.URL{}
	if err := g.decoder.Decode(url); err != nil {
		logger.Log.Sugar().Errorw("Failed to read URL", "error", err)
		return nil, err
	}
	return url, nil
}

func (g *Getter) ReadURLAll() (map[string]models.URL, error) {
	mapData := make(map[string]models.URL, 1)
	idNum := 1
	for g.scanner.Scan() {
		lineData := g.scanner.Bytes()
		url := models.URL{}
		err := json.Unmarshal(lineData, &url)
		if err != nil {
			logger.Log.Sugar().Errorw("Failed Unmarshal URL", "error", err)
		}
		url.ID = uint(idNum)
		mapData[url.ShortURL] = url
		idNum++

	}
	if err := g.scanner.Err(); err != nil {
		logger.Log.Sugar().Error("Failed scan map")
		return nil, err
	}
	return mapData, nil
}

func (g *Getter) Close() error {
	return g.file.Close()
}
