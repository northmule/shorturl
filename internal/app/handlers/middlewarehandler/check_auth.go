package middlewarehandler

import (
	"context"
	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/app/util/user"
	"net/http"
	"time"
)

const defaultUUID = "a4a45d8d-cd8b-47a7-a7a1-4bafcf3d1111"

type CheckAuth struct {
	storage url.StorageInterface
	session *storage.SessionStorage
}

func NewCheckAuth(storage url.StorageInterface, session *storage.SessionStorage) *CheckAuth {
	return &CheckAuth{
		storage: storage,
		session: session,
	}
}

// AuthEveryone выдаст куку с id 1
func (c *CheckAuth) AuthEveryone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet && req.URL.Path == "/api/user/urls" {
			next.ServeHTTP(res, req)
			return
		}
		authorizationToken := req.Header.Get("Authorization")
		var token string
		var userUUID string
		var tokenExp time.Time
		if sessionUserUUID, ok := c.session.Get(authorizationToken); ok {
			userUUID = sessionUserUUID
			tokenExp = time.Now().Add(auntificator.HMACTokenExp)
		} else {
			userUUID = defaultUUID
			token, tokenExp = auntificator.GenerateToken(userUUID, auntificator.HMACTokenExp, auntificator.HMACSecretKey)
			_, err := c.storage.CreateUser(models.User{
				Name:     "test_user",
				UUID:     userUUID,
				Login:    "test_user" + userUUID,
				Password: user.PasswordHash(userUUID),
			})
			if err != nil {
				logger.LogSugar.Errorf("Failed to create user: %v", err)
				return
			}
			c.session.Add(token, userUUID)
		}

		http.SetCookie(res, &http.Cookie{
			Name:    auntificator.CookieAuthName,
			Value:   token,
			Expires: tokenExp,
			Path:    "/",
		})

		res.Header().Set("content-type", "text/plain; charset=utf-8")
		res.Header().Set("Authorization", token)

		ctx := context.WithValue(req.Context(), AppContext.KeyContext, userUUID)

		reqWithContext := req.WithContext(ctx)
		next.ServeHTTP(res, reqWithContext)
	})
}
