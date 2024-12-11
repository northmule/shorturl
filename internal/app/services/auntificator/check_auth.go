package auntificator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage/models"
)

// CheckAuth структура.
type CheckAuth struct {
	userCreator UserCreator
}

// UserCreator интерфейс создания пользователей.
type UserCreator interface {
	CreateUser(user models.User) (int64, error)
}

// NewCheckAuth конструктор
func NewCheckAuth(userCreator UserCreator) *CheckAuth {
	return &CheckAuth{userCreator: userCreator}
}

// ResultCheckAuth результаты работы функции
type ResultCheckAuth struct {
	UserUUID   string
	Token      string
	TokenExp   time.Time
	AuthString string
	IsNewUser  bool
}

// Auth авторизация пользователя.
func (c *CheckAuth) Auth(authorizationToken string) (*ResultCheckAuth, error) {

	res := &ResultCheckAuth{}

	var userUUID, token string
	var exp time.Time
	if authorizationToken == "" {
		userUUID = uuid.NewString()
		token, exp = GenerateToken(userUUID, HMACTokenExp, HMACSecretKey)
		logger.LogSugar.Infof("Cookies have not been transferred, I am creating a new user with a uuid %s", userUUID)
		res.IsNewUser = true
	} else {
		cookieValues := strings.Split(authorizationToken, ":")
		if len(cookieValues) < 2 {
			logger.LogSugar.Infof("The user's UID was not found in the cookie %s", authorizationToken)
			return nil, errors.New("missing user token")
		}
		cookieToken := cookieValues[0]
		userUUID = cookieValues[1]
		logger.LogSugar.Infof("I found cookies for a user with a uuid %s", userUUID)
		if !ValidateToken(userUUID, cookieToken, HMACSecretKey) {
			userUUID = uuid.NewString()
			logger.LogSugar.Infof("The token failed validation for the user with uuid %s. Creating a new user", userUUID)
			token, exp = GenerateToken(userUUID, HMACTokenExp, HMACSecretKey)
			res.IsNewUser = true
		}

	}
	c.createUser(userUUID)

	res.AuthString = fmt.Sprintf("%s:%s", token, userUUID)
	res.Token = token
	res.UserUUID = userUUID
	res.TokenExp = exp

	return res, nil

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
