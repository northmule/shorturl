package middlewarehandler

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/northmule/shorturl/internal/app/handlers/auth"
	"github.com/northmule/shorturl/internal/app/logger"
	"net/http"
)

// роут который авторизует любой запрос
const authRoute = "/api/auth_hmac_everyone"

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == authRoute {
			next.ServeHTTP(res, req)
			return
		}

		cookieAuth, err := req.Cookie(auth.CookieAuthName)
		if err != nil {
			logger.LogSugar.Infof("Ожидалось значение cookie %s", auth.CookieAuthName)
			http.Redirect(res, req, authRoute, http.StatusSeeOther)
			return
		}
		token := cookieAuth.Value
		tokenSign, err := hex.DecodeString(token)
		if err != nil {
			http.Error(res, "Значение cookie не распознано", http.StatusUnauthorized)
			logger.LogSugar.Infof("Не удалось раскодировать token %s", token)

			return
		}
		hashed := hmac.New(sha512.New, []byte(auth.HMACSecretKey))
		hashed.Write([]byte(""))
		expectedSign := hashed.Sum(nil)

		if !hmac.Equal(tokenSign[10:], expectedSign) {
			logger.LogSugar.Infof("Токен %s пользователя подписан другим ключом", token)
			http.Redirect(res, req, authRoute, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(res, req)
	})
}
