package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
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
func (f *FileStorage) Add(url models.URL) (int64, error) {
	modelRaw, err := json.Marshal(url)
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}
	modelJSON := string(modelRaw)

	f.cacheValues = append(f.cacheValues, modelJSON)
	_, err = f.file.WriteString(modelJSON + "\n")
	if err != nil {
		logger.LogSugar.Errorf("Ошибка записи строки %s в файл %s", modelJSON, f.file.Name())
	}

	return 1, nil
}

func (f *FileStorage) CreateUser(user models.User) (int64, error) {
	return 0, nil
}

func (f *FileStorage) LikeURLToUser(urlID int64, userUUID string) error {
	//todo
	return nil
}

func (f *FileStorage) MultiAdd(urls []models.URL) error {
	for _, url := range urls {
		_, err := f.Add(url)
		if err != nil {
			return err
		}
	}
	return nil
}
func (f *FileStorage) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (f *FileStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	for _, value := range f.cacheValues {
		if strings.Contains(value, fmt.Sprintf("\"%s\"", shortURL)) {
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
		if strings.Contains(value, fmt.Sprintf("\"%s\"", url)) {
			url := models.URL{}
			err := json.Unmarshal([]byte(value), &url)
			if err != nil {
				return nil, err
			}
			return &url, nil
		}
	}
	return new(models.URL), nil
}

func (f *FileStorage) FindUserByLoginAndPasswordHash(login string, password string) (*models.User, error) {
	return nil, nil
}

func (f *FileStorage) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}

func (f *FileStorage) Ping() error {
	_, err := os.Stat(f.file.Name())
	if err != nil {
		return err
	}
	return nil
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
