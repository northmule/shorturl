package handlers

import (
	"net/http"
)

func EncodeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "expected get request", http.StatusBadRequest)
		return
	}
	id := req.PathValue("id")
	if id == "" {
		http.Error(res, "expected id value", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusTemporaryRedirect)
	// todo декодировать id и вренуть
	res.Write([]byte("https:ya.ru"))
}
