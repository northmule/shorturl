package storage

import (
	"github.com/northmule/shorturl/internal/app/storage/models"
	"time"
)

// SessionStorage данные по авторизованным пользователям
type SessionStorage struct {
	Values map[string]SessionValue
}

type SessionValue struct {
	User        models.User
	TokenExpiry time.Time
	Token       string
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Values: make(map[string]SessionValue, 100),
	}
}

func (s *SessionValue) IsExpired() bool {
	return s.TokenExpiry.Before(time.Now())
}
