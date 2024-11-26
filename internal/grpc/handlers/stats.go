package handlers

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StatsHandler обработка запросов статистики
type StatsHandler struct {
	contract.UnimplementedStatsHandlerServer
	finderStats handlers.FinderStats
}

// NewStatsHandler конструктор
func NewStatsHandler(finderStats handlers.FinderStats) *StatsHandler {
	instance := &StatsHandler{
		finderStats: finderStats,
	}

	return instance
}

// Stats показывает статистику по пользователям и URL-ам
func (s *StatsHandler) Stats(ctx context.Context, request *empty.Empty) (*contract.StatsResponse, error) {
	var err error

	response := &contract.StatsResponse{}
	response.Users, err = s.finderStats.GetCountUser()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error GetCountUser()")
	}
	response.Urls, err = s.finderStats.GetCountShortURL()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error GetCountShortURL()")
	}

	return response, nil
}
