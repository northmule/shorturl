package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
	_ "go.uber.org/mock/mockgen/model"
)

// CodeErrorDuplicateKey код ошибки с дублем записи в БД.
const CodeErrorDuplicateKey = "23505"

// DBQuery общие методы для работы с хранилищем.
type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// StorageQuery общий интерфес хранилища
type StorageQuery interface {
	// Add добавляет URL.
	Add(url models.URL) (int64, error)
	// CreateUser создание пользователя.
	CreateUser(user models.User) (int64, error)
	// LikeURLToUser связывает пользователя с ссылкой.
	LikeURLToUser(urlID int64, userUUID string) error
	// FindByShortURL поиск по короткой ссылке.
	FindByShortURL(shortURL string) (*models.URL, error)
	// FindByURL поиск по URL.
	FindByURL(url string) (*models.URL, error)
	// Ping проверка соединения с БД.
	Ping() error
	// MultiAdd вставка массива адресов.
	MultiAdd(urls []models.URL) error
	// FindUrlsByUserID поиск ссылок пользователя
	FindUrlsByUserID(userUUID string) (*[]models.URL, error)
	// SoftDeletedShortURL пометка ссылки как удалённой.
	SoftDeletedShortURL(userUUID string, shortURL ...string) error
}

// PostgresStorage хранилище в БД.
type PostgresStorage struct {
	DB    DBQuery
	RawDB *sql.DB
}

// NewPostgresStorage конструктор подключения к БД.
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	// Example: "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	instance := &PostgresStorage{
		DB:    db,
		RawDB: db,
	}

	return instance, err
}

// Add добавление нового значения.
func (p *PostgresStorage) Add(url models.URL) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var urlID int64
	// ON CONFLICT (url) where deleted_at IS NULL DO UPDATE SET url=$2
	err := p.DB.QueryRowContext(ctx, "insert into url_list (short_url, url) values ($1, $2) returning id", url.ShortURL, url.URL).Scan(&urlID)
	return urlID, err
}

// CreateUser добавление нового значения.
func (p *PostgresStorage) CreateUser(user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx, `
			insert into users (name, login, password, uuid) values ($1, $2, $3, $4) ON CONFLICT (uuid) DO UPDATE SET uuid = $4 returning id`, user.Name, user.Login, user.Password, user.UUID)
	return 0, err
}

// LikeURLToUser Связывание URL с пользователем.
func (p *PostgresStorage) LikeURLToUser(urlID int64, userUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx, `insert into user_short_url (user_id, url_id) values ((select id from users where uuid=$1 limit 1), $2)`, userUUID, urlID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
	}
	return err
}

// FindByShortURL поиск по короткой ссылке.
func (p *PostgresStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, short_url, url, deleted_at from url_list where short_url = $1 limit 1",
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
	var deletedAt sql.NullTime
	if rows.Next() {
		err := rows.Scan(&url.ID, &url.ShortURL, &url.URL, &deletedAt)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindByShortURL(%s) произошла ошибка %s", shortURL, err)
			return nil, err
		}
	}
	if deletedAt.Valid {
		url.DeletedAt = deletedAt.Time
	}
	return &url, nil
}

// FindByURL поиск по URL.
func (p *PostgresStorage) FindByURL(url string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, short_url, url from url_list where url = $1 and deleted_at is null limit 1",
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

// Ping проверка соединения.
func (p *PostgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	return p.DB.PingContext(ctx)
}

// MultiAdd Вставка значений в бд пачками.
func (p *PostgresStorage) MultiAdd(urls []models.URL) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	prepareInsert, err := tx.PrepareContext(ctx, `insert into url_list (short_url, url) values ($1, $2) ON CONFLICT (url) where deleted_at IS NULL DO NOTHING;`)
	if err != nil {
		return err
	}
	for _, url := range urls {
		_, err = prepareInsert.ExecContext(ctx, url.ShortURL, url.URL)
		if err != nil {
			logger.LogSugar.Errorf("Значение %#v не добавлено в таблицу url_list", url)
			return errors.Join(err, tx.Rollback())
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// FindUserByLoginAndPasswordHash Поиск пользователя.
func (p *PostgresStorage) FindUserByLoginAndPasswordHash(login string, passwordHash string) (*models.User, error) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, name, login, password from users where login = $1 and password = $2 and deleted_at is null limit 1",
		login,
		passwordHash,
	)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUserByLoginAndPasswordHash(%s) произошла ошибка %s", login, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUserByLoginAndPasswordHash(%s) произошла ошибка %s", login, err)
		return nil, err
	}
	user := models.User{}
	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindUserByLoginAndPasswordHash(%s) произошла ошибка %s", login, err)
			return nil, err
		}
	}

	return &user, nil
}

// FindUrlsByUserID поиск URL-s.
func (p *PostgresStorage) FindUrlsByUserID(userUUID string) (*[]models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		`select ul.id, ul.short_url, ul.url from url_list as ul
				left join user_short_url as usu on usu.url_id=ul.id
				where usu.user_id=(select id from users where uuid=$1 limit 1) order by ul.id asc`,
		userUUID,
	)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUrlsByUserID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUrlsByUserID(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.ShortURL, &url.URL)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindUrlsByUserID(%s) произошла ошибка %s", userUUID, err)
			return nil, err
		}
		urls = append(urls, url)
	}

	return &urls, nil
}

// SoftDeletedShortURL Отметка об удалении ссылки.
func (p *PostgresStorage) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx, `update url_list set deleted_at=now() where short_url = ANY($1)
				and id in (
					select uu.url_id from user_short_url as uu where uu.user_id =
					                                    (select us.id from users as us where us.uuid=$2 limit 1)
	)`, shortURL, userUUID)
	return err
}

// GetCountShortURL кол-во сокращенных URL
func (p *PostgresStorage) GetCountShortURL() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var cnt int64
	rows, err := p.DB.QueryContext(ctx, `select count(*) as cnt from url_list`)
	if err != nil {
		return cnt, err
	}
	err = rows.Err()
	if err != nil {
		return cnt, err
	}
	if rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			return cnt, err
		}
	}

	return cnt, nil
}

// GetCountUser кол-во пользвателей
func (p *PostgresStorage) GetCountUser() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var cnt int64
	rows, err := p.DB.QueryContext(ctx, `select count(*) as cnt from users`)
	if err != nil {
		return cnt, err
	}
	err = rows.Err()
	if err != nil {
		return cnt, err
	}
	if rows.Next() {
		err = rows.Scan(&cnt)
		if err != nil {
			return cnt, err
		}
	}

	return cnt, nil
}
