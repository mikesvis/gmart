package server

import (
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/mikesvis/gmart/internal/handler"
	"github.com/mikesvis/gmart/internal/middleware"
)

func NewRouter(h *handler.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(chiMiddleware.Timeout(60 * time.Second))
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
