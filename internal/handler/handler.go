package handler

import (
	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/service/accural"
	"github.com/mikesvis/gmart/internal/service/order"
	"github.com/mikesvis/gmart/internal/service/user"
)

type Handler struct {
	config  *config.Config
	user    *user.Service
	order   *order.Service
	accural *accural.Service
}

func NewHandler(config *config.Config, user *user.Service, order *order.Service, accural *accural.Service) *Handler {
	return &Handler{config, user, order, accural}
}
