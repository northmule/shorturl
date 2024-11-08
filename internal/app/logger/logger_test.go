package logger

import (
	"testing"
)

func TestInitLogger_InvalidLevel(t *testing.T) {
	err := InitLogger("invalid")
	if err == nil {
		t.Error("InitLogger() error")
	}

}
