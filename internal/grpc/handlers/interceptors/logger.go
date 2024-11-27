package interceptors

import (
	"context"
	"time"

	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

// Logger логгер запросов
type Logger struct {
	l Info
}

// Info интерфейс с методами
type Info interface {
	Infof(template string, args ...interface{})
}

// NewLogger конструктор
func NewLogger(l Info) *Logger {
	return &Logger{l: l}
}

// LogStart начало запроса
func (l *Logger) LogStart(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = utils.AppendMData(ctx, mData.RequestTime, time.Now().String())
	l.l.Infof("Запрос: %s", info.FullMethod)
	l.l.Infof("Req: %v", req)
	l.l.Infof("Ctx: %v", ctx)
	return handler(ctx, req)
}

// LogEnd конец запроса
func (l *Logger) LogEnd(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	l.l.Infof("Обработка запрос: %s завершена", info.FullMethod)

	md, _ := metadata.FromIncomingContext(ctx)
	mdValues := md.Get(mData.RequestTime)
	startTime, _ := time.Parse(time.RFC3339Nano, mdValues[0])
	endTime := time.Since(startTime).String()

	l.l.Infof("Время: %v", endTime)
	return handler(ctx, req)
}
