package handlers

import (
	"net/http"

	"github.com/northmule/shorturl/internal/app/logger"
	"github.com/northmule/shorturl/internal/app/services/url"
)

type PingHandler struct {
	storage url.StorageInterface
}

func NewPingHandler(storage url.StorageInterface) *PingHandler {
	return &PingHandler{
		storage: storage,
	}
}

// CheckStorageConnect обработка запроса проверки соединения с БД /ping
func (p *PingHandler) CheckStorageConnect(res http.ResponseWriter, req *http.Request) {
	err := p.storage.Ping()
	if err != nil {
		http.Error(res, "no connect db", http.StatusInternalServerError)
		logger.LogSugar.Errorf("CheckStorageConnect Не удалось подключиться к БД %s", err)
		return
	}
	res.Write([]byte("Ok"))
	res.WriteHeader(http.StatusOK)
}
