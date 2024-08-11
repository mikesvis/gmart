package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mikesvis/gmart/internal/handler"
)

func NewRouter(h *handler.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middlewares...)

	r.Use(middlewares...)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.UserRegister)
		r.Post("/login", h.UserLogin)
	})

	return r
}
