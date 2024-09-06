package storage

// SessionStorage данные по авторизованным пользователям
type SessionStorage struct {
	Values map[string]string
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Values: make(map[string]string, 100),
	}
}
