package handlers

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/northmule/shorturl/internal/app/handlers"
	"github.com/northmule/shorturl/internal/grpc/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockBadUserFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockBadUserFinder) GetCountShortURL() (int64, error) {
	return 1, nil
}

// GetCountUser кол-во пользвателей
func (s *mockBadUserFinder) GetCountUser() (int64, error) {
	return 0, errors.New("error")
}

type mockBadURLsFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockBadURLsFinder) GetCountShortURL() (int64, error) {
	return 0, errors.New("error")
}

// GetCountUser кол-во пользвателей
func (s *mockBadURLsFinder) GetCountUser() (int64, error) {
	return 1, nil
}

type mockFinder struct {
}

// GetCountShortURL кол-во сокращенных URL
func (s *mockFinder) GetCountShortURL() (int64, error) {
	return 1, nil
}

// GetCountUser кол-во пользвателей
func (s *mockFinder) GetCountUser() (int64, error) {
	return 1, nil
}

func TestStatsHandler_Stats(t *testing.T) {

	tests := []struct {
		name         string
		finder       handlers.FinderStats
		expectedCode codes.Code
	}{
		{
			name:         "error_GetCountUser",
			finder:       new(mockBadUserFinder),
			expectedCode: codes.Internal,
		},
		{
			name:         "error_GetCountShortURL",
			finder:       new(mockBadURLsFinder),
			expectedCode: codes.Internal,
		},
		{
			name:         "ok",
			finder:       new(mockFinder),
			expectedCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpc.NewServer()
			contract.RegisterStatsHandlerServer(s, NewStatsHandler(tt.finder))
			ctx := context.Background()

			dopts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithContextDialer(registerServer(s)),
			}

			conn, err := grpc.NewClient(":///test.server", dopts...)

			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			client := contract.NewStatsHandlerClient(conn)
			r := &empty.Empty{}
			_, err = client.Stats(ctx, r)

			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.expectedCode {
					t.Error("error code: expected", tt.expectedCode, "received", er.Code())
				}
			}
		})
	}
}
