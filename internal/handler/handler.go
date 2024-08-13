package handler

import (
	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/service/order"
	"github.com/mikesvis/gmart/internal/service/user"
)

type Handler struct {
	config *config.Config
	user   *user.Service
	order  *order.Service
}

func NewHandler(config *config.Config, user *user.Service, order *order.Service) *Handler {
	return &Handler{config, user, order}
}
