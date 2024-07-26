package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"io"
	"os"
	"strings"
)

// FileStorage структура хранилища
type FileStorage struct {
	file        *os.File
	scanner     *bufio.Scanner
	cacheValues []string
}

func NewFileStorage(file *os.File) *FileStorage {
	instance := &FileStorage{
		file:        file,
		scanner:     bufio.NewScanner(file),
		cacheValues: make([]string, 0),
	}
	instance.restoreStorage()
	return instance
}

// Add добавление нового значения
func (f *FileStorage) Add(url models.URL) error {
	modelRaw, err := json.Marshal(url)
	if err != nil {
		logger.LogSugar.Error(err)
		return err
	}
	modelJSON := string(modelRaw)

	if len(f.cacheValues) > 0 {
		_, err = f.file.Seek(0, io.SeekStart)
		if err != nil {
			logger.LogSugar.Errorf("Ошибка усечения файла")
			return err
		}
		err = f.file.Truncate(0)
		if err != nil {
			logger.LogSugar.Errorf("Ошибка усечения файла")
			return err
		}
	}

	if len(f.cacheValues) == 0 {
		f.cacheValues = append(f.cacheValues, modelJSON)
	}

	for _, value := range f.cacheValues {
		_, err = f.file.WriteString(value + "\n")
		if err != nil {
			logger.LogSugar.Errorf("Ошибка записи строки %s в файл %s", value, f.file.Name())
		}

		if !strings.Contains(value, modelJSON) {
			f.cacheValues = append(f.cacheValues, modelJSON)
			_, err = f.file.WriteString(modelJSON + "\n")
			if err != nil {
				logger.LogSugar.Errorf("Ошибка записи строки %s в файл %s", modelJSON, f.file.Name())
			}
		}
	}

	return nil
}

// FindByShortURL поиск по короткой ссылке
func (f *FileStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	for _, value := range f.cacheValues {
		if strings.Contains(value, shortURL) {
			url := models.URL{}
			err := json.Unmarshal([]byte(value), &url)
			if err != nil {
				logger.LogSugar.Errorf("Ошибка json.Unmarshal: %s", value)
				return nil, err
			}
			return &url, nil
		}
	}

	return nil, fmt.Errorf("the short link was not found")
}

// FindByURL поиск по URL
func (f *FileStorage) FindByURL(url string) (*models.URL, error) {
	for _, value := range f.cacheValues {
		if strings.Contains(value, url) {
			url := models.URL{}
			err := json.Unmarshal([]byte(value), &url)
			if err != nil {
				return nil, err
			}
			return &url, nil
		}
	}
	return nil, fmt.Errorf("the url link was not found")
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}

// restoreStorage восстановит бд из переданного значения
func (f *FileStorage) restoreStorage() {
	for f.scanner.Scan() {
		lineData := string(f.scanner.Bytes())
		f.cacheValues = append(f.cacheValues, lineData)
	}
	if err := f.scanner.Err(); err != nil {
		logger.LogSugar.Errorf("При восстановлении храналица, обнаружены ошибки: %s", err)
	}
}
