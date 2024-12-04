package interceptors

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/northmule/shorturl/internal/grpc/handlers/metadata"
	"github.com/northmule/shorturl/internal/grpc/handlers/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckAuth структура.
type CheckAuth struct {
	userCreator                       middlewarehandler.UserCreator
	session                           storage.SessionAdapter
	checkAuthExpectedMethods          []string
	accessVerificationExpectedMethods []string
}

// NewCheckAuth конструктор структуры.
func NewCheckAuth(userCreator middlewarehandler.UserCreator, session storage.SessionAdapter) *CheckAuth {
	return &CheckAuth{
		userCreator:                       userCreator,
		session:                           session,
		checkAuthExpectedMethods:          []string{"/contract.ShortenerHandler/Shortener", "/contract.ShortenerHandler/ShortenerJSON", "/contract.UserUrlsHandler/View", "/contract.UserUrlsHandler/Delete"},
		accessVerificationExpectedMethods: []string{"/contract.UserUrlsHandler/View"},
	}
}

// AuthEveryone авторизация пользователя.
func (c *CheckAuth) AuthEveryone(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if !isMethodExpected(info, c.checkAuthExpectedMethods) {
		return handler(ctx, req)
	}

	authorizationToken := utils.GetUserToken(ctx)

	var userUUID string
	if authorizationToken == "" {
		userUUID = uuid.NewString()
		token, _ := auntificator.GenerateToken(userUUID, auntificator.HMACTokenExp, auntificator.HMACSecretKey)
		logger.LogSugar.Infof("Cookies have not been transferred, I am creating a new user with a uuid %s", userUUID)
		c.createUser(userUUID)

		ctx = utils.AppendMData(ctx, metadata.Authorization, fmt.Sprintf("%s:%s", token, userUUID))
	} else {
		cookieValues := strings.Split(authorizationToken, ":")
		if len(cookieValues) < 2 {
			logger.LogSugar.Infof("The user's UID was not found in the cookie %s", authorizationToken)
			return nil, status.Error(codes.Unauthenticated, "missing user token")
		}
		cookieToken := cookieValues[0]
		userUUID = cookieValues[1]
		logger.LogSugar.Infof("I found cookies for a user with a uuid %s", userUUID)
		if !auntificator.ValidateToken(userUUID, cookieToken, auntificator.HMACSecretKey) {
			userUUID = uuid.NewString()
			logger.LogSugar.Infof("The token failed validation for the user with uuid %s. Creating a new user", userUUID)
			token, _ := auntificator.GenerateToken(userUUID, auntificator.HMACTokenExp, auntificator.HMACSecretKey)

			ctx = utils.AppendMData(ctx, metadata.Authorization, fmt.Sprintf("%s:%s", token, userUUID))
		}

		c.createUser(userUUID)
	}

	ctx = utils.AppendMData(ctx, metadata.UserUUID, userUUID)

	return handler(ctx, req)

}

// AccessVerificationUserUrls проверка доступа пользователя.
func (c *CheckAuth) AccessVerificationUserUrls(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if !isMethodExpected(info, c.accessVerificationExpectedMethods) {
		return handler(ctx, req)
	}

	authorizationToken := utils.GetUserToken(ctx)

	if authorizationToken == "" {
		logger.LogSugar.Infof("The user's UUID was not found in the cookie %s when requesting /api/user/urls", authorizationToken)
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	return handler(ctx, req)
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

func isMethodExpected(info *grpc.UnaryServerInfo, expected []string) bool {
	for _, method := range expected {
		if method == info.FullMethod {
			return true
		}
	}

	return false
}
