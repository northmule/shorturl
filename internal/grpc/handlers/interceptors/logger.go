package interceptors

import (
	"context"
	"time"

	"github.com/northmule/shorturl/internal/app/logger"
	mData "github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

// Logger логгер запросов
type Logger struct {
}

// NewLogger конструктор
func NewLogger() *Logger {
	return &Logger{}
}

// LogStart начало запроса
func (l *Logger) LogStart(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = utils.AppendMData(ctx, mData.RequestTime, time.Now().String())
	logger.LogSugar.Infof("Запрос: %s", info.FullMethod)
	logger.LogSugar.Infof("Req: %v", req)
	logger.LogSugar.Infof("Ctx: %v", ctx)
	return handler(ctx, req)
}

// LogEnd конец запроса
func (l *Logger) LogEnd(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.LogSugar.Infof("Обработка запрос: %s завершена", info.FullMethod)

	md, _ := metadata.FromIncomingContext(ctx)
	mdValues := md.Get(mData.RequestTime)
	startTime, _ := time.Parse(time.RFC3339Nano, mdValues[0])
	endTime := time.Since(startTime).String()

	logger.LogSugar.Infof("Время: %v", endTime)
	return handler(ctx, req)
}
