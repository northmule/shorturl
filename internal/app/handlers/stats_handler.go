package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/northmule/shorturl/internal/app/logger"
)

// StatsHandler обработка запросов статистики
type StatsHandler struct {
	finderStats StatsFinder
}

// StatsFinder интерфейс поиска данных
type StatsFinder interface {
	GetCountShortURL() (int64, error)
	GetCountUser() (int64, error)
}

// NewStatsHandler конструктор
func NewStatsHandler(finderStats StatsFinder) *StatsHandler {
	instance := &StatsHandler{
		finderStats: finderStats,
	}

	return instance
}

// ResponseViewStats ответ хэндлера
type ResponseViewStats struct {
	Urls  int64 `json:"urls"`
	Users int64 `json:"users"`
}

// ViewStats показывает статистику по пользователям и URL-ам
func (s *StatsHandler) ViewStats(res http.ResponseWriter, req *http.Request) {
	var err error
	var responseView ResponseViewStats
	responseView.Users, err = s.finderStats.GetCountUser()
	if err != nil {
		logger.LogSugar.Error("error GetCountUser()")
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}
	responseView.Urls, err = s.finderStats.GetCountShortURL()
	if err != nil {
		logger.LogSugar.Error("error GetCountShortURL()")
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(responseView)
	if err != nil {
		logger.LogSugar.Error("error json marshal response")
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseBytes)
	if err != nil {
		logger.LogSugar.Error("error write data")
		http.Error(res, "error write", http.StatusInternalServerError)
		return
	}
}
