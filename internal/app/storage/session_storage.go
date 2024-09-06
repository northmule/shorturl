package storage

import "sync"

// SessionStorage данные по авторизованным пользователям
type SessionStorage struct {
	Values map[string]string
	mx     sync.RWMutex
}

func NewSessionStorage() Session {
	return &SessionStorage{
		Values: make(map[string]string, 100),
	}
}

type Session interface {
	Add(key string, value string)
	Get(key string) (string, bool)
	GetAll() map[string]string
}

func (s *SessionStorage) Add(key string, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.Values[key] = value
}

func (s *SessionStorage) Get(key string) (string, bool) {
	sessionUserUUID, ok := s.Values[key]
	return sessionUserUUID, ok
}

func (s *SessionStorage) GetAll() map[string]string {
	return s.Values
}
