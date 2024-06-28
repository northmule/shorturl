package models

// Url Модель данных
type Url struct {
	Id       uint   `gorm:"primary_key;auto_increment" json:"id"`
	ShortUrl string `gorm:"type:varchar(255);unique_index" json:"short_url"`
	Url      string `gorm:"type:varchar(255);unique_index" json:"url"`
}
