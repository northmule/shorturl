package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
	_ "go.uber.org/mock/mockgen/model"
	"time"
)

type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
}

type PostgresStorage struct {
	DB DBQuery
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

	return instance, instance.createTable()
}

// Add добавление нового значения
func (p *PostgresStorage) Add(url models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx, "insert into url_list (short_url, url) values ($1, $2)", url.ShortURL, url.URL)
	if err != nil {
		logger.LogSugar.Errorf("Значение %#v не добавлено в таблицу url_list", url)
		return err
	}
	return nil
}

// FindByShortURL поиск по короткой ссылке
func (p *PostgresStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, short_url, url from url_list where short_url = $1 limit 1",
		shortURL,
	)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindByShortURL(%s) произошла ошибка %s", shortURL, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindByShortURL(%s) произошла ошибка %s", shortURL, err)
		return nil, err
	}
	url := models.URL{}
	if rows.Next() {
		err := rows.Scan(&url.ID, &url.ShortURL, &url.URL)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindByShortURL(%s) произошла ошибка %s", shortURL, err)
			return nil, err
		}
	}

	return &url, nil
}

// FindByURL поиск по URL
func (p *PostgresStorage) FindByURL(url string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, short_url, url from url_list where url = $1 limit 1",
		url,
	)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindByURL(%s) произошла ошибка %s", url, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindByShortURL(%s) произошла ошибка %s", url, err)
		return nil, err
	}
	modelURL := models.URL{}
	if rows.Next() {
		err := rows.Scan(&modelURL.ID, &modelURL.ShortURL, &modelURL.URL)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindByURL(%s) произошла ошибка %s", url, err)
			return nil, err
		}
	}

	return &modelURL, nil
}

// Ping проверка соединения
func (p *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}

// MultiAdd Вставка значений в бд пачками
func (p *PostgresStorage) MultiAdd(urls []models.URL) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	prepareInsert, err := tx.PrepareContext(ctx, `insert into url_list (short_url, url) values ($1, $2)`)
	if err != nil {
		return err
	}
	for _, url := range urls {
		_, err := prepareInsert.ExecContext(ctx, url.ShortURL, url.URL)
		if err != nil {
			logger.LogSugar.Errorf("Значение %#v не добавлено в таблицу url_list", url)
			errR := tx.Rollback()
			if errR != nil {
				logger.LogSugar.Errorf("откат транзакции вызвал сбой: %s", errR)
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// createTable создаёт необходимую таблицу при её отсутсвии
func (p *PostgresStorage) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS public.url_list (
					id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
					short_url varchar(100) NOT NULL,
					url varchar(2000) NOT NULL,
					created_at timestamp DEFAULT now() NOT NULL,
					deleted_at timestamp NULL,
					CONSTRAINT url_pk PRIMARY KEY (id)
				);
					CREATE INDEX IF NOT EXISTS url_short_url_idx ON public.url_list USING btree (short_url);
					CREATE INDEX IF NOT EXISTS url_url_idx ON public.url_list USING btree (url)`,
	)
	if err != nil {
		logger.LogSugar.Errorf("Ошибка создания базы данных: %s", err)
		return err
	}
	logger.LogSugar.Info("Создание таблицы завершено")
	return nil
}
