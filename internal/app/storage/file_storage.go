package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// FileStorage файловое хранилище.
type FileStorage struct {
	file        *os.File
	scanner     *bufio.Scanner
	cacheValues []string
	users       *os.File
	deletedURLs *os.File
}

// NewFileStorage конструктор хранилища.
func NewFileStorage(file *os.File) *FileStorage {
	instance := &FileStorage{
		file:        file,
		scanner:     bufio.NewScanner(file),
		cacheValues: make([]string, 0),
	}

	usersFileName := file.Name() + "user.json"

	fileUsers, err := os.OpenFile(usersFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.LogSugar.Errorf("Failed to open file %s: error: %s", usersFileName, err)
		return nil
	}
	deletedURLsFileName := file.Name() + "deleted-urls.json"
	fileDeletedURLs, err := os.OpenFile(deletedURLsFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.LogSugar.Errorf("Failed to open file %s: error: %s", deletedURLsFileName, err)
		return nil
	}
	instance.users = fileUsers
	instance.deletedURLs = fileDeletedURLs
	instance.restoreStorage()
	return instance
}

// Add добавление нового значения.
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

// CreateUser создает пользователя.
func (f *FileStorage) CreateUser(user models.User) (int64, error) {
	modelRaw, err := json.Marshal(user)
	if err != nil {
		logger.LogSugar.Error(err)
		return 0, err
	}
	modelJSON := string(modelRaw)

	_, err = f.users.WriteString(modelJSON + "\n")
	if err != nil {
		logger.LogSugar.Errorf("Ошибка записи строки %s в файл %s", modelJSON, f.users.Name())
	}
	return 0, nil
}

// LikeURLToUser Связывание URL с пользователем.
func (f *FileStorage) LikeURLToUser(urlID int64, userUUID string) error {
	//todo
	return nil
}

// MultiAdd Вставка массива.
func (f *FileStorage) MultiAdd(urls []models.URL) error {
	for _, url := range urls {
		_, err := f.Add(url)
		if err != nil {
			return err
		}
	}
	return nil
}

// SoftDeletedShortURL Отметка об удалении ссылки.
func (f *FileStorage) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	return nil
}

// FindByShortURL поиск по короткой ссылке.
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

// FindByURL поиск по URL.
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

// FindUserByLoginAndPasswordHash Поиск пользователя.
func (f *FileStorage) FindUserByLoginAndPasswordHash(login string, password string) (*models.User, error) {
	return nil, nil
}

// FindUrlsByUserID поиск URL-s.
func (f *FileStorage) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	return nil, nil
}

// Close закрытие файла
func (f *FileStorage) Close() error {
	return f.file.Close()
}

// Ping проверка доступности.
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

// GetCountShortURL кол-во сокращенных URL
func (f *FileStorage) GetCountShortURL() (int64, error) {
	return int64(len(f.cacheValues)), nil
}

// GetCountUser кол-во пользвателей
func (f *FileStorage) GetCountUser() (int64, error) {
	userFile, err := os.Open(f.users.Name())
	if err != nil {
		return 0, err
	}
	b := bufio.NewScanner(userFile)
	var cnt int64
	for b.Scan() {
		cnt++
	}

	return cnt, nil
}
