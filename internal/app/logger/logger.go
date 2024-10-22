package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// LogSugar Глобальный логгер.
var LogSugar *Logger

// Logger Логгер для реализации в запросах
type Logger struct {
	*zap.SugaredLogger
}

// LogEntry Логгер для реализации в запросах
type LogEntry struct {
	*zap.SugaredLogger
}

// NewLogger конструктор.
func NewLogger(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	LogSugar = &Logger{appLogger.Sugar()}
	return LogSugar, nil
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
