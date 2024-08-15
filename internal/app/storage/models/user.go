package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Urls     []URL  `json:"urls"`
}
