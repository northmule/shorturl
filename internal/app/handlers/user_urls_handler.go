package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/northmule/shorturl/internal/app/handlers/auth"
	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
	"github.com/northmule/shorturl/internal/app/storage"
	"github.com/northmule/shorturl/internal/app/storage/models"
	"net/http"
)

type UserURLsHandler struct {
	storage url.StorageInterface
	session *storage.SessionStorage
}

func NewUserUrlsHandler(storage url.StorageInterface, sessionStorage *storage.SessionStorage) *UserURLsHandler {
	instance := UserURLsHandler{
		storage: storage,
		session: sessionStorage,
	}
	return &instance
}

type ResponseView struct {
	ShortUrl    string `json:"short_Url"`
	OriginalUrl string `json:"original_url"`
}

// View коротки ссылки пользователя
func (u *UserURLsHandler) View(res http.ResponseWriter, req *http.Request) {
	cookieAuth, err := req.Cookie(auth.CookieAuthName)
	if err != nil {
		logger.LogSugar.Infof("Ожидалось значение cookie %s", auth.CookieAuthName)
		return
	}
	token := cookieAuth.Value
	//tokenSign, err := hex.DecodeString(token)
	//if err != nil {
	//	http.Error(res, "Значение cookie не распознано", http.StatusUnauthorized)
	//	logger.LogSugar.Infof("Не удалось раскодировать token %s", token)
	//	return
	//}
	//userId := string(tokenSign[:10]) // при создании токена первые 10 байт отводятся для id (см Auth)
	//beforeHashData := make([]byte, 10)
	//copy(beforeHashData, userId)
	//hashed := hmac.New(sha512.New, []byte(auth.HMACSecretKey))
	//hashed.Write(beforeHashData)
	//expectedSign := hashed.Sum(beforeHashData)
	//
	//if !bytes.Equal(tokenSign, expectedSign) {
	//	logger.LogSugar.Infof("Токен %s пользователя подписан другим ключом", token)
	//	return
	//}

	userId, err := GetUserIdFromCookie(cookieAuth)

	var session storage.SessionValue
	for _, value := range u.session.Values {
		if token == value.Token {
			session = value
			break
		}
	}

	if session.User.ID == 0 {
		logger.LogSugar.Infof("Пользователь с токеном %s не найден", token)
		userById, err := u.storage.FindUserById(userId)
		if err != nil || userById.ID == 0 {
			http.Error(res, "Пользователь с таким id не найден", http.StatusUnauthorized)
			logger.LogSugar.Infof("Пользователь с таким id: %d не найден", userId)
			return
		}
		session.User = *userById

	}
	if session.IsExpired() {
		http.Error(res, "время жизни токена не действительно", http.StatusUnauthorized)
		logger.LogSugar.Infof("Пользователь найден, но время токена истекло: %s для пользователя %s", session.TokenExpiry, session.User.Login)
		return
	}
	err = u.fillURLs(&session.User)
	if err != nil {
		http.Error(res, "Ошибка получения ссылок пользователя", http.StatusInternalServerError)
		logger.LogSugar.Error(err)
		return
	}
	res.Header().Set("content-type", "application/json")

	if len(session.User.Urls) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var responseList []ResponseView
	for _, urlItem := range session.User.Urls {
		responseList = append(responseList, ResponseView{
			ShortUrl:    urlItem.ShortURL,
			OriginalUrl: urlItem.URL,
		})
	}
	responseURLs, err := json.Marshal(responseList)
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseURLs)
	if err != nil {
		http.Error(res, "error write data", http.StatusInternalServerError)
		return
	}
}

func (u *UserURLsHandler) fillURLs(user *models.User) error {
	if user.Urls != nil {
		return nil
	}
	urls, err := u.storage.FindUrlsByUserId(user.ID)
	if err != nil {
		return err
	}
	user.Urls = *urls
	return nil
}

func GetUserIdFromCookie(c *http.Cookie) (int, error) {
	token := c.Value
	tokenSign, err := hex.DecodeString(token)
	if err != nil {
		logger.LogSugar.Infof("Не удалось раскодировать token %s", token)
		return 0, err
	}
	hashed := hmac.New(sha512.New, []byte(auth.HMACSecretKey))
	hashed.Write([]byte(""))
	expectedSign := hashed.Sum(nil)

	if !hmac.Equal(tokenSign[10:], expectedSign) {
		logger.LogSugar.Infof("Токен %s пользователя подписан другим ключом", token)
		return 0, fmt.Errorf("токен %s пользователя подписан другим ключом", token)
	}
	numberUserId := binary.BigEndian.Uint64(tokenSign[:10]) // при создании токена первые 10 байт отводятся для id (см Auth)
	return int(numberUserId), nil
}
