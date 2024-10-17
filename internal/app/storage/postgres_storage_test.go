package storage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/northmule/shorturl/internal/app/logger"
	mocks "github.com/northmule/shorturl/internal/app/storage/mock"
	"github.com/northmule/shorturl/internal/app/storage/models"
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
				fmt.Println("Recovered in f", r)
			}
		}()
		storage.Add(models.URL{})
	})

}

func TestPostgresStorage_FindByShortURL(t *testing.T) {
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
				fmt.Println("Recovered in f", r)
			}
		}()
		storage.FindByShortURL("")
	})
}

func TestPostgresStorage_FindByURL(t *testing.T) {
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
				fmt.Println("Recovered in f", r)
			}
		}()
		storage.FindByURL("")
	})
}

func TestPostgresStorage_Ping(t *testing.T) {
	t.Run("Вызов_Ping", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		m.EXPECT().PingContext(gomock.Any()).Return(nil)
		storage := PostgresStorage{DB: m}
		storage.Ping()
	})
}

func TestPostgresStorage_createTable(t *testing.T) {
	logger.NewLogger("info")
	t.Run("Создание_таблиц", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		result := &testResult{}
		m.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil).Times(3)
		tx := new(sql.Tx)
		m.EXPECT().Begin().Return(tx, nil)
		storage := PostgresStorage{DB: m}
		defer func() {
			// вызов Commit
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		storage.createTable()
	})
}

// При использовании github.com/DATA-DOG/go-sqlmock
func TestPostgresStorage_createTable_by_sqlmock(t *testing.T) {
	logger.NewLogger("info")
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	result := sqlmock.NewResult(1, 1)
	mock.ExpectBegin()
	// При сравнении полного запроса будет ошибка, по этому достаточно части запроса из оригинальной функции.
	//См vendor/github.com/DATA-DOG/go-sqlmock/query.go:45
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS public.url_list").WillReturnResult(result)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS public.users").WillReturnResult(result)
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS public.user_short_url").WillReturnResult(result)
	mock.ExpectCommit()

	storage := PostgresStorage{DB: db}
	err = storage.createTable()
	if err != nil {
		t.Errorf("an error '%s' was not expected when creating the table", err)
	}

}
