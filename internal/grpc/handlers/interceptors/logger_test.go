package interceptors

import (
	"context"
	"testing"
	"time"

	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type mockInfo struct {
	mock.Mock
	data []string
}

func (m *mockInfo) Infof(template string, args ...interface{}) {
	m.data = append(m.data, template)
}

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Invoke(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	return args.Get(0), args.Error(1)
}

func TestLogStart(t *testing.T) {
	mockLogger := new(mockInfo)
	mockHandler := new(mockHandler)

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/method"}

	mockHandler.On("Invoke", mock.Anything, req).Return(nil, nil)

	l := NewLogger(mockLogger)

	_, _ = l.LogStart(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return mockHandler.Invoke(ctx, req)
	})
	assert.Equal(t, 3, len(mockLogger.data))
}

func TestLogEnd(t *testing.T) {
	mockLogger := new(mockInfo)
	mockHandler := new(mockHandler)

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{FullMethod: "/test/method"}

	startTime := time.Now().Format(time.RFC3339Nano)
	ctx = utils.AppendMData(ctx, mData.RequestTime, startTime)

	mockHandler.On("Invoke", mock.Anything, req).Return(nil, nil)

	l := NewLogger(mockLogger)
	_, _ = l.LogEnd(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return mockHandler.Invoke(ctx, req)
	})

	assert.Equal(t, 2, len(mockLogger.data))
}
