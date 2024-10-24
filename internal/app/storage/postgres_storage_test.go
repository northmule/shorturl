package storage

import (
	"context"
	"database/sql"
	"testing"

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
		storage.Ping()
	})
}
