package auntificator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/northmule/shorturl/internal/app/logger"
)

// Общие константы сервиса
const (
	// CookieAuthName название куки авторизации.
	CookieAuthName = "shorturl_session"
	// HMACTokenExp время жизни токена.
	HMACTokenExp = time.Hour * 600
	// HMACSecretKey секретный ключ генерации токена.
	HMACSecretKey = "super_secret_key"
)

// GetUserToken получить токен из запроса.
func GetUserToken(req *http.Request) string {
	token := req.Header.Get("Authorization")
	if token == "" {
		cookieAuth, err := req.Cookie(CookieAuthName)
		if err != nil {
			logger.LogSugar.Infof("Ожидалось значение cookie %s", CookieAuthName)
			return ""
		}
		token = cookieAuth.Value
	}

	return token
}

// GenerateToken сгенерировать новый токен по параметрам.
func GenerateToken(userUUID string, exp time.Duration, secretKey string) (string, time.Time) {
	hashed := hmac.New(sha256.New, []byte(secretKey))
	hashed.Write([]byte(userUUID))
	token := hex.EncodeToString(hashed.Sum(nil))
	tokenExp := time.Now().Add(exp)
	return token, tokenExp
}

// ValidateToken проверка валидности токена.
func ValidateToken(userUUID string, token string, secretKey string) bool {
	logger.LogSugar.Infof("Проверка токена %s для пользователя %s", token, userUUID)
	tokenSign, err := hex.DecodeString(token)
	if err != nil {
		logger.LogSugar.Infof("Не удалось раскодировать token %s", token)
		return false
	}
	hashed := hmac.New(sha256.New, []byte(secretKey))
	hashed.Write([]byte(userUUID))
	expectedSign := hashed.Sum(nil)

	if !hmac.Equal(tokenSign, expectedSign) {
		logger.LogSugar.Infof("Токен %s пользователя подписан другим ключом", token)
		return false
	}
	logger.LogSugar.Infof("Токен %s прошёл проверку", token)
	return true
}
