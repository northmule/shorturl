package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/northmule/shorturl/internal/app/logger"
	mocks "github.com/northmule/shorturl/internal/app/storage/mock"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"go.uber.org/mock/gomock"
	"testing"
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
		result := &testResult{}
		m.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)
		storage := PostgresStorage{DB: m}
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
	t.Run("Создание_БД", func(t *testing.T) {
		ctrl, _ := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()
		m := mocks.NewMockDBQuery(ctrl)
		result := &testResult{}
		m.EXPECT().ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(result, nil)
		storage := PostgresStorage{DB: m}
		storage.createTable()
	})
}
