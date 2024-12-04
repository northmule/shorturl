package middlewarehandler

import (
	"net/http"

	"github.com/northmule/shorturl/config"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
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
		trustedService := auntificator.NewTrustedSubnet(c.configApp)
		err = trustedService.GrantAccess(req.Header.Get("X-Real-IP"))
		if err != nil {
			logger.LogSugar.Warn(err.Error())
			res.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(res, req)
	})
}
