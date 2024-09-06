package auntificator

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/northmule/shorturl/internal/app/logger"
	"net/http"
	"time"
)

// idCookieSize размер в байтах места под id пользователя
const idCookieSize = 4
const CookieAuthName = "shorturl_session"
const HMACTokenExp = time.Hour * 600
const HMACSecretKey = "super_secret_key"

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

func GenerateToken(userUUID string, exp time.Duration, secretKey string) (string, time.Time) {
	beforeHashData := make([]byte, 0, idCookieSize)
	copy(beforeHashData, userUUID)
	hashed := hmac.New(sha512.New, []byte(secretKey))
	hashed.Write([]byte("rrr"))
	token := hex.EncodeToString(hashed.Sum(beforeHashData))
	tokenExp := time.Now().Add(exp)
	return token, tokenExp
}
