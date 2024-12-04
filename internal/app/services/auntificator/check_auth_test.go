package auntificator

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCreator struct {
	mock.Mock
}

func (m *MockUserCreator) CreateUser(user models.User) (int64, error) {
	args := m.Called(user)
	return 1, args.Error(1)
}

func TestAuthWithEmptyToken(t *testing.T) {
	mockUserCreator := new(MockUserCreator)
	checkAuth := NewCheckAuth(mockUserCreator)

	mockUserCreator.On("CreateUser", mock.Anything).Return(int64(1), nil)

	result, err := checkAuth.Auth("")
	assert.NoError(t, err)
	assert.True(t, result.IsNewUser)
	assert.NotEmpty(t, result.UserUUID)
	assert.NotEmpty(t, result.Token)
	assert.NotEqual(t, time.Time{}, result.TokenExp)

	mockUserCreator.AssertCalled(t, "CreateUser", mock.MatchedBy(func(user models.User) bool {
		return user.UUID == result.UserUUID && user.Name == "test_user" && user.Login == "test_user"+result.UserUUID && user.Password == "password"
	}))
}

func TestAuthWithInvalidToken(t *testing.T) {
	mockUserCreator := new(MockUserCreator)
	checkAuth := NewCheckAuth(mockUserCreator)

	userUUID := uuid.NewString()
	invalidToken := "invalid_token"
	authString := fmt.Sprintf("%s:%s", invalidToken, userUUID)

	mockUserCreator.On("CreateUser", mock.Anything).Return(int64(1), nil)

	result, err := checkAuth.Auth(authString)
	assert.NoError(t, err)
	assert.True(t, result.IsNewUser)
	assert.NotEqual(t, userUUID, result.UserUUID)
	assert.NotEmpty(t, result.Token)
	assert.NotEqual(t, time.Time{}, result.TokenExp)

	mockUserCreator.AssertCalled(t, "CreateUser", mock.MatchedBy(func(user models.User) bool {
		return user.UUID == result.UserUUID && user.Name == "test_user" && user.Login == "test_user"+result.UserUUID && user.Password == "password"
	}))
}
