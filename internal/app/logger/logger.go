package logger

import (
	"go.uber.org/zap"
)

// Log глобальный логгер.
var Log *zap.Logger = zap.NewNop()

// LogSugar глобальный логгер sugar.
var LogSugar *zap.SugaredLogger

// NewLogger конструктор.
func NewLogger(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = appLogger
	LogSugar = Log.Sugar()
	return nil
}
