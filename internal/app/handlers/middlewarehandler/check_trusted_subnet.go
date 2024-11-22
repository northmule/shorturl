package middlewarehandler

import (
	"net"
	"net/http"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
)

// CheckTrustedSubnet проверка запроса на принадлежность к сети
type CheckTrustedSubnet struct {
	configApp *config.Config
}

// NewCheckTrustedSubnet конструктор
func NewCheckTrustedSubnet(configApp *config.Config) *CheckTrustedSubnet {
	return &CheckTrustedSubnet{
		configApp: configApp,
	}
}

// GrantAccess предоставить доступ
func (c *CheckTrustedSubnet) GrantAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var err error
		if c.configApp.TrustedSubnet == "" {
			logger.LogSugar.Infof("доверенная сеть не заданна, доступ ограничен")
			res.WriteHeader(http.StatusForbidden)
			return
		}

		var expectedIP net.IP
		var actualIP net.IP
		var expectedNet *net.IPNet

		actualIP = net.ParseIP(req.Header.Get("X-Real-IP"))
		if actualIP == nil {
			logger.LogSugar.Infof("не передан IP адрес, доступ ограничен")
			res.WriteHeader(http.StatusForbidden)
			return
		}

		expectedIP, expectedNet, err = net.ParseCIDR(c.configApp.TrustedSubnet)
		if err != nil {
			logger.LogSugar.Infof("адрес конфигурации не распознан, доступ ограничен")
			res.WriteHeader(http.StatusForbidden)
			return
		}

		if ok := expectedIP.Equal(actualIP); ok {
			next.ServeHTTP(res, req)
		}

		if ok := expectedNet.Contains(actualIP); !ok {
			logger.LogSugar.Infof("адрес не является разрешённым, доступ ограничен")
			res.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(res, req)
	})
}
