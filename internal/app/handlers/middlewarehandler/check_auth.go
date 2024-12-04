package middlewarehandler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// CheckAuth структура.
type CheckAuth struct {
	userCreator UserCreator
	session     storage.SessionAdapter
}

// NewCheckAuth конструктор структуры.
func NewCheckAuth(userCreator UserCreator, session storage.SessionAdapter) *CheckAuth {
	return &CheckAuth{
		userCreator: userCreator,
		session:     session,
	}
}

// UserCreator интерфейс создания пользователей.
type UserCreator interface {
	CreateUser(user models.User) (int64, error)
}

// AccessVerificationUserUrls проверка доступа пользователя.
func (c *CheckAuth) AccessVerificationUserUrls(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authorizationToken := auntificator.GetUserToken(req)

		if authorizationToken == "" {
			logger.LogSugar.Infof("UUID пользователя в куке не найден %s при запросе /api/user/urls", authorizationToken)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}

// AuthEveryone авторизация пользователя.
func (c *CheckAuth) AuthEveryone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		checkAuthService := auntificator.NewCheckAuth(c.userCreator)

		authorizationToken := auntificator.GetUserToken(req)
		authResult, err := checkAuthService.Auth(authorizationToken)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		if authResult.IsNewUser {
			res = c.authorization(res, authResult.UserUUID, authResult.Token, authResult.TokenExp)
		}

		res.Header().Set("content-type", "text/plain; charset=utf-8")
		ctx := context.WithValue(req.Context(), AppContext.KeyContext, authResult.UserUUID)

		reqWithContext := req.WithContext(ctx)
		next.ServeHTTP(res, reqWithContext)
	})
}

func (c *CheckAuth) authorization(res http.ResponseWriter, userUUID string, token string, tokenExp time.Time) http.ResponseWriter {
	tokenValue := fmt.Sprintf("%s:%s", token, userUUID)
	http.SetCookie(res, &http.Cookie{
		Name:    auntificator.CookieAuthName,
		Value:   tokenValue,
		Expires: tokenExp,
		Secure:  false,
		Path:    "/",
	})
	res.Header().Set("Authorization", tokenValue)
	return res
}

func (c *CheckAuth) createUser(userUUID string) {
	_, err := c.userCreator.CreateUser(models.User{
		Name:     "test_user",
		UUID:     userUUID,
		Login:    "test_user" + userUUID,
		Password: "password",
	})
	if err != nil {
		logger.LogSugar.Errorf("Failed to create user: %v", err)
		return
	}
}
