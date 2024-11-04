package middlewarehandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	AppContext "github.com/northmule/shorturl/internal/app/context"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"github.com/stretchr/testify/mock"
)

type MockUserCreator struct {
	mock.Mock
	user models.User
}

func (m *MockUserCreator) CreateUser(user models.User) (int64, error) {
	m.On("CreateUser", user).Return(int64(1), nil)
	m.Called(user)
	return 1, nil
}

func TestAuthEveryone_NewUser(t *testing.T) {
	_ = logger.InitLogger("fatal")
	userCreator := new(MockUserCreator)
	session := storage.NewSessionStorage()
	handler := NewCheckAuth(userCreator, session)

	req, err := http.NewRequest("GET", "/api/user/urls", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	authHandler := handler.AuthEveryone(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(AppContext.KeyContext).(string)
		if userUUID == "" {
			t.Error("Expected user UUID in context, but got empty string")
		}
		w.WriteHeader(http.StatusOK)
	}))

	authHandler.ServeHTTP(res, req)

	userCreator.AssertCalled(t, "CreateUser", mock.AnythingOfType("models.User"))

	if res.Header().Get("Authorization") == "" {
		t.Error("Expected Authorization header to be set, but got empty string")
	}
}

func TestAuthEveryone_TokenNoValid(t *testing.T) {
	_ = logger.InitLogger("fatal")
	userCreator := new(MockUserCreator)
	session := storage.NewSessionStorage()
	handler := NewCheckAuth(userCreator, session)

	req, err := http.NewRequest("GET", "/api/user/urls", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "sfsdfdsfsdf:1111111-222222-33333-444444")
	res := httptest.NewRecorder()

	authHandler := handler.AuthEveryone(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value(AppContext.KeyContext).(string)
		if userUUID == "" {
			t.Error("Expected user UUID in context, but got empty string")
		}
		w.WriteHeader(http.StatusOK)
	}))

	authHandler.ServeHTTP(res, req)

	userCreator.AssertCalled(t, "CreateUser", mock.AnythingOfType("models.User"))

	if res.Header().Get("Authorization") == "" {
		t.Error("Expected Authorization header to be set, but got empty string")
	}
}
