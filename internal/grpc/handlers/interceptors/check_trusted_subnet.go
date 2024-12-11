package interceptors

import (
	"context"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CheckTrustedSubnet проверка запроса на принадлежность к сети
type CheckTrustedSubnet struct {
	configApp                  *config.Config
	grantAccessExpectedMethods []string
}

// NewCheckTrustedSubnet конструктор
func NewCheckTrustedSubnet(configApp *config.Config) *CheckTrustedSubnet {
	return &CheckTrustedSubnet{
		configApp:                  configApp,
		grantAccessExpectedMethods: []string{"/contract.StatsHandler/Stats"},
	}
}

// GrantAccess предоставить доступ
func (c *CheckTrustedSubnet) GrantAccess(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if !isMethodExpected(info, c.grantAccessExpectedMethods) {
		return handler(ctx, req)
	}

	var err error

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	mdValues := md.Get("X-Real-IP")

	if len(mdValues) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	trustedService := auntificator.NewTrustedSubnet(c.configApp)
	err = trustedService.GrantAccess(mdValues[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "no access")
	}

	return handler(ctx, req)

}
