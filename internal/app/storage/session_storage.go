package storage

import "sync"

// SessionStorage данные по авторизованным пользователям.
type SessionStorage struct {
	Values map[string]string
	mx     sync.RWMutex
}

// NewSessionStorage конструктор.
func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Values: make(map[string]string, 100),
	}
}

// Session метод работы с хранилищем.
type Session interface {
	// Add добавить новую запись.
	Add(key string, value string)
	// Get получить запись по ключу.
	Get(key string) (string, bool)
	// GetAll получить все записи из хранилища.
	GetAll() map[string]string
}

// Add добавить новую запись.
func (s *SessionStorage) Add(key string, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Values[key] = value
}

// Get получить запись по ключу.
func (s *SessionStorage) Get(key string) (string, bool) {
	sessionUserUUID, ok := s.Values[key]
	return sessionUserUUID, ok
}

// GetAll получить все записи из хранилища.
func (s *SessionStorage) GetAll() map[string]string {
	return s.Values
}
