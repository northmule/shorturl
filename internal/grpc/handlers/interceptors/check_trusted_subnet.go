package interceptors

import (
	"context"
	"net"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
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
	if c.configApp.TrustedSubnet == "" {
		logger.LogSugar.Infof("доверенная сеть не заданна, доступ ограничен")
		return nil, status.Error(codes.Unauthenticated, "missing TrustedSubnet")
	}

	var expectedIP net.IP
	var actualIP net.IP
	var expectedNet *net.IPNet

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	mdValues := md.Get("X-Real-IP")

	if len(mdValues) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	actualIP = net.ParseIP(mdValues[0])
	if actualIP == nil {
		logger.LogSugar.Infof("не передан IP адрес, доступ ограничен")
		return nil, status.Error(codes.Unauthenticated, "missing X-Real-IP")
	}

	expectedIP, expectedNet, err = net.ParseCIDR(c.configApp.TrustedSubnet)
	if err != nil {
		logger.LogSugar.Infof("адрес конфигурации не распознан, доступ ограничен")
		return nil, status.Error(codes.Unauthenticated, "missing CIDR")
	}

	if ok := expectedIP.Equal(actualIP); ok {
		return handler(ctx, req)
	}

	if ok := expectedNet.Contains(actualIP); !ok {
		logger.LogSugar.Infof("адрес не является разрешённым, доступ ограничен")
		return nil, status.Error(codes.Unauthenticated, "missing expected IP")
	}

	return handler(ctx, req)

}
