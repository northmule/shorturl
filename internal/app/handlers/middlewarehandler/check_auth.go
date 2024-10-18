package middlewarehandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

const defaultUUID = "a4a45d8d-cd8b-47a7-a7a1-4bafcf3d1111"

type CheckAuth struct {
	storage url.StorageInterface
	session *storage.Session
}

func NewCheckAuth(storage url.StorageInterface, session *storage.Session) *CheckAuth {
	return &CheckAuth{
		storage: storage,
		session: session,
	}
}

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

// AuthEveryone выдаст куку
func (c *CheckAuth) AuthEveryone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authorizationToken := auntificator.GetUserToken(req)

		var userUUID string
		if authorizationToken == "" {

			userUUID = uuid.NewString()
			token, tokenExp := auntificator.GenerateToken(userUUID, auntificator.HMACTokenExp, auntificator.HMACSecretKey)
			logger.LogSugar.Infof("Куки не переданы, создаю нового пользователя с uuid %s", userUUID)
			c.createUser(userUUID)
			res = c.authorization(res, userUUID, token, tokenExp)
		} else {
			cookieValues := strings.Split(authorizationToken, ":")
			if len(cookieValues) < 2 {
				logger.LogSugar.Infof("UUID пользователя в куке не найден %s", authorizationToken)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
			cookieToken := cookieValues[0]
			userUUID = cookieValues[1]
			logger.LogSugar.Infof("Нашёл куки для пользователя с uuid %s", userUUID)
			if !auntificator.ValidateToken(userUUID, cookieToken, auntificator.HMACSecretKey) {
				userUUID = uuid.NewString()
				logger.LogSugar.Infof("Токен не прошёл валидацию для пользователя с uuid %s. Создаю нового пользователя", userUUID)
				token, tokenExp := auntificator.GenerateToken(userUUID, auntificator.HMACTokenExp, auntificator.HMACSecretKey)
				res = c.authorization(res, userUUID, token, tokenExp)
			}

			c.createUser(userUUID)
		}

		res.Header().Set("content-type", "text/plain; charset=utf-8")
		ctx := context.WithValue(req.Context(), AppContext.KeyContext, userUUID)

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
	_, err := c.storage.CreateUser(models.User{
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
