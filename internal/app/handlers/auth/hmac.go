package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	passwordUtil "github.com/northmule/shorturl/internal/app/util/user"
	"net/http"
	"strconv"
	"time"
)

const HMACTokenExp = time.Hour * 600
const HMACSecretKey = "Надо, Федя, надо!"
const CookieAuthName = "shorturl_session"

type HMACAuth struct {
	storage url.StorageInterface
	session *storage.SessionStorage
}

// sessionData данные по авторизованным пользователям
type sessionData struct {
	user        models.User
	tokenExpiry time.Time
	token       string
}

type HMACRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type errorCookieAuthNotFound struct {
	err error
}

func (e *errorCookieAuthNotFound) Error() string {
	return fmt.Sprintf("Cookie %s пустая или не содержит id пользователя: %s", CookieAuthName, e.err)
}

func NewHMACHandler(storage url.StorageInterface, sessionStorage *storage.SessionStorage) *HMACAuth {
	instance := HMACAuth{
		storage: storage,
		session: sessionStorage,
	}
	return &instance
}

// Auth аунтифицирует по логину и паролю
func (h *HMACAuth) Auth(res http.ResponseWriter, req *http.Request) {
	authRequest := HMACRequest{
		Login:    req.FormValue("login"),
		Password: req.FormValue("password"),
	}
	if authRequest.Login == "" || authRequest.Password == "" {
		http.Error(res, "пустые параметры запроса", http.StatusBadRequest)
		logger.LogSugar.Error("Пустые параметры запроса")
		return
	}
	passwordHash := passwordUtil.PasswordHash(authRequest.Password)
	if _, ok := h.session.Values[authRequest.Login+passwordHash]; ok {
		http.Error(res, "Пользователь уже авторизован", http.StatusConflict)
		logger.LogSugar.Infof("Пользователь уже авторизован: %s", authRequest.Login)
		return
	}
	logger.LogSugar.Infof("Данные авторизации: Логин:%s Хэш:%s", authRequest.Login, passwordHash)
	user, err := h.storage.FindUserByLoginAndPasswordHash(authRequest.Login, passwordHash)
	if err != nil {
		http.Error(res, "пользователь не найден", http.StatusNotFound)
		logger.LogSugar.Errorf("Пользователь не найден: %s c хэш: %s не найден", authRequest.Login, passwordHash)
		return
	}

	hashed := hmac.New(sha512.New, []byte(HMACSecretKey))
	// токен с ид пользвоателя
	userId := strconv.Itoa(user.ID)
	beforeHashData := make([]byte, 10)
	copy(beforeHashData, userId)
	hashed.Write([]byte(""))
	token := fmt.Sprintf("%x", hashed.Sum(beforeHashData))
	tokenExp := time.Now().Add(HMACTokenExp)

	h.session.Values[authRequest.Login+passwordHash] = storage.SessionValue{
		User:        *user,
		TokenExpiry: tokenExp,
		Token:       token,
	}

	// Куки должны ставиться до заголовков !!!
	http.SetCookie(res, &http.Cookie{
		Name:    CookieAuthName,
		Value:   token,
		Expires: tokenExp,
		Path:    "/",
	})

	res.Header().Set("content-type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)

	_, err = res.Write([]byte("Вы авторизованы"))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}

// AuthEveryone выдаст куку с id 1
func (h *HMACAuth) AuthEveryone(res http.ResponseWriter, req *http.Request) {
	user := models.User{
		ID:       1,
		Login:    "test",
		Password: "super_password",
	}

	authRequest := HMACRequest{
		Login:    user.Login,
		Password: user.Password,
	}

	passwordHash := passwordUtil.PasswordHash(authRequest.Password)
	if _, ok := h.session.Values[authRequest.Login+passwordHash]; ok {
		logger.LogSugar.Infof("Пользователь уже авторизован: %s", authRequest.Login)
		return
	}
	logger.LogSugar.Infof("Данные авторизации: Логин:%s Хэш:%s", authRequest.Login, passwordHash)

	// токен с ид пользвоателя
	userId := strconv.Itoa(user.ID)
	beforeHashData := make([]byte, 10)
	copy(beforeHashData, userId)
	hashed := hmac.New(sha512.New, []byte(HMACSecretKey))
	hashed.Write([]byte("байты_терабайты"))
	token := fmt.Sprintf("%x", hashed.Sum(beforeHashData))
	tokenExp := time.Now().Add(HMACTokenExp)

	h.session.Values[authRequest.Login+passwordHash] = storage.SessionValue{
		User:        user,
		TokenExpiry: tokenExp,
		Token:       token,
	}

	http.SetCookie(res, &http.Cookie{
		Name:    CookieAuthName,
		Value:   token,
		Expires: tokenExp,
		Path:    "/",
	})

	res.Header().Set("content-type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)

	_, err := res.Write([]byte("Вы авторизованы"))
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}

func (h *HMACAuth) Valid(req *http.Request) (int, error) {
	return 0, nil // todo
}
