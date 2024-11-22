package storage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/shorturl/internal/app/logger"
	mocks "github.com/northmule/shorturl/internal/app/storage/mocks"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type testResult struct {
}

func (t *testResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (t *testResult) RowsAffected() (int64, error) {
	return 0, nil
}

func TestPostgresStorage_Add(t *testing.T) {
	_ = logger.InitLogger("fatal")
	t.Run("Добавление_нового_значения", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		row := &sql.Row{}
		m.EXPECT().QueryRowContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(row)
		storage := PostgresStorage{DB: m}
		defer func() {
			// вызов Next в Add
			if r := recover(); r != nil {
				logger.LogSugar.Infof("Recovered in %v", r)
			}
		}()
		storage.Add(models.URL{})
	})

}

func TestPostgresStorage_FindByShortURL(t *testing.T) {
	_ = logger.InitLogger("fatal")
	t.Run("Поиск_по_короткому_url", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		var rows sql.Rows
		m.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&rows, nil)
		storage := PostgresStorage{DB: m}
		defer func() {
			// вызов Next в FindByShortURL
			if r := recover(); r != nil {
				logger.LogSugar.Infof("Recovered in %v", r)
			}
		}()
		storage.FindByShortURL("")
	})
}

func TestPostgresStorage_FindByURL(t *testing.T) {
	_ = logger.InitLogger("fatal")
	t.Run("Поиск_по_url", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		var rows sql.Rows
		m.EXPECT().QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&rows, nil)
		storage := PostgresStorage{DB: m}
		defer func() {
			// вызов Next в FindByURL
			if r := recover(); r != nil {
				logger.LogSugar.Infof("Recovered in %v", r)
			}
		}()
		storage.FindByURL("")
	})
}

func TestPostgresStorage_Ping(t *testing.T) {
	_ = logger.InitLogger("fatal")
	t.Run("Вызов_Ping", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		m.EXPECT().PingContext(gomock.Any()).Return(nil)
		storage := PostgresStorage{DB: m}
		err := storage.Ping()
		if err != nil {
			fmt.Println(err)
		}
	})
}

type PostgresStorageTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	DB   *sql.DB
	pg   PostgresStorage
}

func (o *PostgresStorageTestSuite) SetupTest() {
	var err error
	o.DB, o.mock, err = sqlmock.New()
	o.pg = PostgresStorage{
		DB:    o.DB,
		RawDB: o.DB,
	}
	require.NoError(o.T(), err)
	_ = logger.InitLogger("fatal")
}

func TestPostgresStorageTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresStorageTestSuite))
}

func (o *PostgresStorageTestSuite) TestMultiAdd() {
	urls := []models.URL{
		{URL: "https://ya.ru/1", ShortURL: "abc123"},
		{URL: "https://ya.ru/2", ShortURL: "abc321"},
	}

	o.mock.ExpectBegin()

	exp := o.mock.ExpectPrepare("insert into")
	for _, url := range urls {
		exp.ExpectExec().WithArgs(url.ShortURL, url.URL).
			WillReturnResult(sqlmock.NewResult(1, 1))

	}
	o.mock.ExpectCommit()
	err := o.pg.MultiAdd(urls)
	require.NoError(o.T(), err)

}

func (o *PostgresStorageTestSuite) TestFindUserByLoginAndPasswordHash() {
	login := "cat"
	pwd := "has_has"

	o.mock.ExpectQuery("select id, name, login, password from users").
		WithArgs(login, pwd).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "login", "password"}).
			AddRow("1", "Кот в Сапогах", "cat", "has_has"))

	user, err := o.pg.FindUserByLoginAndPasswordHash(login, pwd)
	require.NoError(o.T(), err)
	require.Equal(o.T(), login, user.Login)
}

func (o *PostgresStorageTestSuite) TestCreateUser() {
	testUser := models.User{
		Name:     "Кот в Сапогах",
		Login:    "cat",
		Password: "has_has",
		UUID:     "111-222-333",
	}

	o.mock.ExpectExec("insert into users").
		WithArgs(testUser.Name, testUser.Login, testUser.Password, testUser.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := o.pg.CreateUser(testUser)
	require.NoError(o.T(), err)

}

func (o *PostgresStorageTestSuite) TestFindUrlsByUserID() {
	userUUID := "1111-2222-3333-4444"
	o.mock.ExpectQuery("select ul.id, ul.short_url, ul.url").
		WithArgs(userUUID).
		WillReturnRows(sqlmock.NewRows([]string{"ul.id", "ul.short_url", "ul.url"}).
			AddRow("1", "short123", "https://yandex.ru"))
	urls, err := o.pg.FindUrlsByUserID(userUUID)
	require.NoError(o.T(), err)
	require.Equal(o.T(), 1, len(*urls))
}

func (o *PostgresStorageTestSuite) TestGetCountUser() {
	o.mock.ExpectQuery("select count").
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).
			AddRow(43))
	cnt, err := o.pg.GetCountUser()
	require.NoError(o.T(), err)
	require.Equal(o.T(), int64(43), cnt)
}

func (o *PostgresStorageTestSuite) TestGetGetCountShortURL() {
	o.mock.ExpectQuery("select count").
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).
			AddRow(51))
	cnt, err := o.pg.GetCountShortURL()
	require.NoError(o.T(), err)
	require.Equal(o.T(), int64(51), cnt)
}
