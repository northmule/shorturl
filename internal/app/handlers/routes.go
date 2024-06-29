package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// AppRoutes маршруты приложения
func AppRoutes() chi.Router {
	r := chi.NewRouter()
	// Код ошибки по умолчанию, согласно ТЗ
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("method not expect\n"))
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", DecodeHandler)
		r.Get("/{id}", EncodeHandler)
	})

	return r
}
