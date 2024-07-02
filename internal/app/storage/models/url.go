package models

// URL Модель данных
type URL struct {
	ID       uint   `gorm:"primary_key;auto_increment" json:"id"`
	ShortURL string `gorm:"type:varchar(255);unique_index" json:"short_url"`
	URL      string `gorm:"type:varchar(500);unique_index" json:"url"`
}
