package interceptors

import (
	"context"

	"github.com/northmule/shorturl/internal/app/handlers/middlewarehandler"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/auntificator"
	"github.com/northmule/shorturl/internal/app/storage"
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

	checkAuthService := auntificator.NewCheckAuth(c.userCreator)

	authorizationToken := utils.GetUserToken(ctx)

	authResult, err := checkAuthService.Auth(authorizationToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "missing user token")
	}
	if authResult.IsNewUser {
		ctx = utils.AppendMData(ctx, metadata.Authorization, authResult.AuthString)
	}

	ctx = utils.AppendMData(ctx, metadata.UserUUID, authResult.UserUUID)

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

func isMethodExpected(info *grpc.UnaryServerInfo, expected []string) bool {
	for _, method := range expected {
		if method == info.FullMethod {
			return true
		}
	}

	return false
}
