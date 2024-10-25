package handlers

import (
	"net/http"

	"github.com/northmule/shorturl/internal/app/logger"
)

// PingHandler хэндлер для обработки ping запроса.
type PingHandler struct {
	pinger Pinger
}

// NewPingHandler конструктор.
func NewPingHandler(pinger Pinger) *PingHandler {
	return &PingHandler{
		pinger: pinger,
	}
}

// Pinger Интерфейс првоерки соединения.
type Pinger interface {
	// Ping проверка соединения с БД.
	Ping() error
}

// CheckStorageConnect обработка запроса проверки соединения с БД /ping.
// @Summary Проверка подключения к БД
// @Success 200 {json} {ok}
// @Failure 500 {string} string "Не удалось подключиться к БД"
// @Router /ping [get]
func (p *PingHandler) CheckStorageConnect(res http.ResponseWriter, req *http.Request) {
	err := p.pinger.Ping()
	if err != nil {
		http.Error(res, "no connect db", http.StatusInternalServerError)
		logger.LogSugar.Errorf("CheckStorageConnect Не удалось подключиться к БД %s", err)
		return
	}
	res.Write([]byte("Ok"))
	res.WriteHeader(http.StatusOK)
}
