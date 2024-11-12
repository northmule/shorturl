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

func TestAccessVerificationUserUrls(t *testing.T) {

	_ = logger.InitLogger("fatal")
	userCreator := new(MockUserCreator)
	session := storage.NewSessionStorage()
	checkAuth := NewCheckAuth(userCreator, session)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Токен есть
	t.Run("with_token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req.Header.Add("Cookie", "gophermart_session=20106c00221e4fe38582e36a12327fdfc05b904bbad8d30d238a8f8323fcf90d:fbbad27c-16b3-48e3-a455-785074e45981; shorturl_session=2e41a78f9851029e40e85e70ac38d24ca06dcc8bcbb1da0e3619f17edd6c050a:431300ac-c58c-4dcf-941c-e47ca43511ba")

		rr := httptest.NewRecorder()
		handler := checkAuth.AccessVerificationUserUrls(nextHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	// Токена нет
	t.Run("without_token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		rr := httptest.NewRecorder()
		handler := checkAuth.AccessVerificationUserUrls(nextHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}

	})
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

func TestCreateUser(t *testing.T) {
	userCreator := new(MockUserCreator)
	session := storage.NewSessionStorage()
	handler := NewCheckAuth(userCreator, session)
	handler.createUser("1111-2222-33333-44444")
	userCreator.AssertCalled(t, "CreateUser", mock.AnythingOfType("models.User"))
}
