package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PingHandler хэндлер для обработки ping запроса.
type PingHandler struct {
	contract.UnimplementedPingHandlerServer
	pinger handlers.Pinger
}

// NewPingHandler конструктор.
func NewPingHandler(pinger handlers.Pinger) *PingHandler {
	return &PingHandler{pinger: pinger}
}

// CheckStorageConnect обработка запроса проверки соединения с БД .
func (p *PingHandler) CheckStorageConnect(ctx context.Context, request *empty.Empty) (*contract.CheckStorageConnectResponse, error) {
	err := p.pinger.Ping()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "no connect db")
	}
	response := contract.CheckStorageConnectResponse{}
	response.Ok = true

	return &response, nil
}
