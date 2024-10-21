package models

// User модель пользователя
type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	UUID     string `json:"uuid"`
	Urls     []URL  `json:"urls"`
}
