package auth

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	passwordUtil "github.com/northmule/shorturl/internal/app/util/user"
	"io"
	"net/http"
	"time"
)

const JWTTokenExp = time.Hour * 3
const JWTSecretKey = "super_secret_key"

type JWTAuth struct {
	storage url.StorageInterface
}

// Claims утверждение
type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

type JWTRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

func NewJWTHandler(storage url.StorageInterface) *JWTAuth {
	instance := JWTAuth{
		storage: storage,
	}
	return &instance
}

// Auth аунтифицирует по логину и паролю
func (j *JWTAuth) Auth(res http.ResponseWriter, req *http.Request) {
	bodyValue, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "error read bodyValue", http.StatusBadRequest)
		logger.LogSugar.Error("Тепло запроса не прочитанно: ", err)
		return
	}

	defer req.Body.Close()

	var authRequest JWTRequest
	if err = json.Unmarshal(bodyValue, &authRequest); err != nil {
		http.Error(res, "error unmarshal json request", http.StatusBadRequest)
		logger.LogSugar.Error("Не корректное тело запроса авторизации: ", bodyValue)
		return
	}
	if authRequest.Login == "" || authRequest.Password == "" {
		http.Error(res, "пустые параметры запроса", http.StatusBadRequest)
		logger.LogSugar.Error("Пустые параметры запроса")
		return
	}
	passwordHash := passwordUtil.PasswordHash(authRequest.Password)
	logger.LogSugar.Infof("Данные авторизации: Логин:%s Хэш:%s", authRequest.Login, passwordHash)
	user, err := j.storage.FindUserByLoginAndPasswordHash(authRequest.Login, passwordHash)
	if err != nil {
		http.Error(res, "пользователь не найден", http.StatusNotFound)
		logger.LogSugar.Errorf("Пользователь не найден: %s c хэш: %s не найден", authRequest.Login, passwordHash)
		return
	}
	// Токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTTokenExp)),
		},
		UserID: user.ID,
	})

	// подпись токена ключом
	tokenValue, err := token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	responseString, err := json.Marshal(JWTResponse{
		Token: tokenValue,
	})
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseString)
	if err != nil {
		http.Error(res, "error write data", http.StatusBadRequest)
		return
	}
}
