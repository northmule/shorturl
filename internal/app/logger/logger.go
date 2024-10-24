package logger

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// LogSugar Глобальный логгер.
var (
	LogSugar *Logger
	once     sync.Once
)

// Logger Логгер для реализации в запросах
type Logger struct {
	*zap.SugaredLogger
}

// LogEntry Логгер для реализации в запросах
type LogEntry struct {
	*zap.SugaredLogger
}

// InitLogger конструктор.
func InitLogger(level string) error {
	var isError bool
	var err error

	once.Do(
		func() {
			lvl, err := zap.ParseAtomicLevel(level)
			if err != nil {
				isError = true
				return
			}
			cfg := zap.NewDevelopmentConfig()
			cfg.Level = lvl
			appLogger, err := cfg.Build()
			if err != nil {
				isError = true
				return
			}
			LogSugar = &Logger{appLogger.Sugar()}
		})

	if isError {
		return err
	}

	return nil
}

// NewLogEntry Конструктор
func (l *Logger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &LogEntry{
		l.SugaredLogger,
	}
}

// Print Печать
func (l *Logger) Print(v ...interface{}) {
	l.Info(v...)
}

// Write Печать сообщения
func (l *LogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Infof("Информация о запросе: Статус: %d. Байт: %d. Заголовки: %#v. Время: %d. Дополнительно: %#v", status, bytes, header, elapsed, extra)
}

// Panic Сообщение при панике
func (l *LogEntry) Panic(v interface{}, stack []byte) {
	l.Infof("Паника: %#v. Трейс: %s", v, string(stack))
}
