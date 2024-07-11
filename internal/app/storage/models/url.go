package models

// URL Модель данных
type URL struct {
	ID       uint   `json:"id,omitempty"`
	ShortURL string `json:"short_url"`
	URL      string `json:"url"`
}
