package handler

import (
	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/service/user"
)

type Handler struct {
	config *config.Config
	user   *user.Service
}

func NewHandler(config *config.Config, user *user.Service) *Handler {
	return &Handler{config, user}
}
