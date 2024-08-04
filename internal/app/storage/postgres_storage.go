package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"time"
)

type PostgresStorage struct {
	DB *sql.DB
}

// NewPostgresStorage PostgresStorage настройка подключения к БД
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	// Example: "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	instance := &PostgresStorage{
		DB: db,
	}
	return instance, nil
}

// Add добавление нового значения
func (p *PostgresStorage) Add(url models.URL) error {
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (p *PostgresStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	return nil, nil
}

// FindByURL поиск по URL
func (p *PostgresStorage) FindByURL(url string) (*models.URL, error) {
	return nil, nil
}

func (p *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}
