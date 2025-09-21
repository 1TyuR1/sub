package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterPing(r chi.Router) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})
}
