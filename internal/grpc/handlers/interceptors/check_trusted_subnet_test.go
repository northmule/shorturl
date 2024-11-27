package interceptors

import (
	"context"
	"testing"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockTrustedHandler struct {
	mock.Mock
}

func (m *mockTrustedHandler) Invoke(ctx context.Context, req interface{}) (interface{}, error) {
	return nil, nil
}

func TestGrantAccess(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		code    codes.Code
		message string
		ctx     func() context.Context
	}{
		{
			name: "NoTrustedSubnet",
			cfg: config.Config{
				TrustedSubnet: "",
			},
			code:    codes.Unauthenticated,
			message: "missing TrustedSubnet",
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "NoTrustedSubnet",
			cfg: config.Config{
				TrustedSubnet: "invalid-cidr",
			},
			code:    codes.Unauthenticated,
			message: "missing metadata",
			ctx: func() context.Context {
				return context.Background()
			},
		},
		{
			name: "IPInSubnet",
			cfg: config.Config{
				TrustedSubnet: "192.168.1.0/24",
			},
			code:    codes.OK,
			message: "",
			ctx: func() context.Context {
				md := metadata.New(map[string]string{"X-Real-IP": "192.168.1.199"})
				ctx := metadata.NewIncomingContext(context.Background(), md)
				return ctx
			},
		},
	}
	mockHandler := new(mockTrustedHandler)
	_ = logger.InitLogger("fatal")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inter := NewCheckTrustedSubnet(&tt.cfg)

			ctx := tt.ctx()
			req := "test request"
			info := &grpc.UnaryServerInfo{FullMethod: "/contract.StatsHandler/Stats"}

			_, err := inter.GrantAccess(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
				return mockHandler.Invoke(ctx, req)
			})
			if er, ok := status.FromError(err); ok {
				if er.Code() != tt.code {
					t.Error("error code: expected", tt.code, "received", er.Code())
				}
				if er.Message() != tt.message {
					t.Error("error message: expected", tt.message, "received", er.Message())
				}
			}
		})
	}

}
