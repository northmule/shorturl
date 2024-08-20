package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/migrations"
	"github.com/northmule/shorturl/internal/app/storage/models"
	_ "go.uber.org/mock/mockgen/model"
	"time"
)

const CodeErrorDuplicateKey = "23505"

type DBQuery interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
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
func (p *PostgresStorage) Add(url models.URL) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var urlID int64
	// ON CONFLICT (url) where deleted_at IS NULL DO UPDATE SET url=$2
	err := p.DB.QueryRowContext(ctx, "insert into url_list (short_url, url) values ($1, $2) returning id", url.ShortURL, url.URL).Scan(&urlID)
	return urlID, err
}

// CreateUser добавление нового значения
func (p *PostgresStorage) CreateUser(user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	var insertID int64
	_ = p.DB.QueryRowContext(
		ctx,
		"insert into users (name, login, password, uuid) values ($1, $2, $3, $4) returning id",
		user.Name,
		user.Login,
		user.Password,
		user.UUID,
	).Scan(&insertID)
	return insertID, nil
}

func (p *PostgresStorage) LikeURLToUser(urlID int64, userUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	_, err := p.DB.ExecContext(ctx, "insert into user_short_url (user_id, url_id) values ((select id from users where uuid=$1 limit 1), $2)", userUUID, urlID)
	return err
}

// FindByShortURL поиск по короткой ссылке
func (p *PostgresStorage) FindByShortURL(shortURL string) (*models.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		"select id, short_url, url from url_list where short_url = $1 and deleted_at is null limit 1",
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

	prepareInsert, err := tx.PrepareContext(ctx, `insert into url_list (short_url, url) values ($1, $2) ON CONFLICT (url) where deleted_at IS NULL DO NOTHING;`)
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

func (p *PostgresStorage) FindUserByLoginAndPasswordHash(login string, passwordHash string) (*models.User, error) {
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
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindUserByLoginAndPasswordHash(%s) произошла ошибка %s", login, err)
			return nil, err
		}
	}

	return &user, nil
}

func (p *PostgresStorage) FindUrlsByUserId(userUUID string) (*[]models.URL, error) {
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
		logger.LogSugar.Errorf("При вызове FindUrlsByUserId(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUrlsByUserId(%s) произошла ошибка %s", userUUID, err)
		return nil, err
	}
	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.ShortURL, &url.URL)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindUrlsByUserId(%s) произошла ошибка %s", userUUID, err)
			return nil, err
		}
		urls = append(urls, url)
	}

	return &urls, nil
}

func (p *PostgresStorage) FindUserById(userId int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	rows, err := p.DB.QueryContext(
		ctx,
		`select id, name, login, password from users where id = $1 and deleted_at is null limit 1`,
		userId,
	)
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUserById(%d) произошла ошибка %s", userId, err)
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		logger.LogSugar.Errorf("При вызове FindUserById(%d) произошла ошибка %s", userId, err)
		return nil, err
	}
	var user models.User
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password)
		if err != nil {
			logger.LogSugar.Errorf("При обработке значений в FindUserById(%d) произошла ошибка %s", userId, err)
			return nil, err
		}
	}

	return &user, nil
}

// createTable создаёт необходимую таблицу при её отсутсвии
func (p *PostgresStorage) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()

	logger.LogSugar.Info("Попытка создать таблицу url_list")
	_, err := p.DB.ExecContext(ctx, migrations.Migrations01)
	if err != nil {
		logger.LogSugar.Errorf("Ошибка создания таблицы: %s", err)
		return err
	}

	logger.LogSugar.Info("Попытка создать таблицу users")
	_, err = p.DB.ExecContext(ctx, migrations.Migrations02)
	if err != nil {
		logger.LogSugar.Errorf("Ошибка создания таблицы: %s", err)
		return err
	}

	logger.LogSugar.Info("Попытка создать таблицу user_short_url")
	_, err = p.DB.ExecContext(ctx, migrations.Migrations03)
	if err != nil {
		logger.LogSugar.Errorf("Ошибка создания таблицы: %s", err)
		return err
	}

	logger.LogSugar.Info("Создание таблиц завершено")
	return nil
}
