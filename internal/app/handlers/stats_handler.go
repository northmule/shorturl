package handlers

import (
	"encoding/json"
	"net/http"
)

// StatsHandler обработка запросов статистики
type StatsHandler struct {
	finderStats FinderStats
}

// FinderStats интерфейс поиска данных
type FinderStats interface {
	GetCountShortURL() (int64, error)
	GetCountUser() (int64, error)
}

// NewStatsHandler конструктор
func NewStatsHandler(finderStats FinderStats) *StatsHandler {
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
		http.Error(res, "error GetCountUser()", http.StatusInternalServerError)
		return
	}
	responseView.Urls, err = s.finderStats.GetCountShortURL()
	if err != nil {
		http.Error(res, "error GetCountShortURL()", http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(responseView)
	if err != nil {
		http.Error(res, "error json marshal response", http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseBytes)
	if err != nil {
		http.Error(res, "error write data", http.StatusInternalServerError)
		return
	}
}
