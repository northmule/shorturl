package workers

import (
	"testing"
	"time"

	"github.com/northmule/shorturl/internal/app/logger"
)

// Mock Deleter interface
type MockDeleter struct {
	DeleteCalled bool
}

func (m *MockDeleter) SoftDeletedShortURL(userUUID string, shortURL ...string) error {
	m.DeleteCalled = true
	return nil
}

func TestNewWorker(t *testing.T) {
	mockDeleter := &MockDeleter{}
	stopChan := make(chan struct{})

	worker := NewWorker(mockDeleter, stopChan)

	if worker.deleter != mockDeleter {
		t.Errorf("Expected deleter to be %v, but got %v", mockDeleter, worker.deleter)
	}

	if worker.jobChan == nil {
		t.Errorf("Expected jobChan to be initialized, but it is nil")
	}

	if worker.stopChan != stopChan {
		t.Errorf("Expected stopChan to be %v, but got %v", stopChan, worker.stopChan)
	}

}

func TestWorker_Del(t *testing.T) {
	err := logger.InitLogger("fatal")
	if err != nil {
		t.Error(err)
	}
	mockDeleter := &MockDeleter{}

	stopChan := make(chan struct{})

	worker := NewWorker(mockDeleter, stopChan)
	worker.Del("user1", []string{"url1", "url2"})
	time.Sleep(100 * time.Millisecond)

	if !mockDeleter.DeleteCalled {
		t.Error("Expected mockDeleter.Delete to be called")
	}
}
