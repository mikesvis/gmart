package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mikesvis/gmart/internal/handler"
	"github.com/mikesvis/gmart/internal/middleware"
)

func NewRouter(h *handler.Handler, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(middlewares...)

	r.Use(middlewares...)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", h.UserRegister)
		r.Post("/login", h.UserLogin)
		r.With(middleware.Auth).Post("/orders", h.CreateUserOrder)
		r.With(middleware.Auth).Get("/orders", h.GetUserOrders)
		r.With(middleware.Auth).Get("/balance", h.GetUserBalance)
		r.With(middleware.Auth).Post("/balance/withdraw", h.WithdrawForOrder)
		r.With(middleware.Auth).Get("/withdrawals", h.GetUserWithdrawals)
	})

	return r
}
