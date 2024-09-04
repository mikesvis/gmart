package handler

import (
	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/service/accrual"
	"github.com/mikesvis/gmart/internal/service/order"
	"github.com/mikesvis/gmart/internal/service/user"
	"go.uber.org/zap"
)

type Handler struct {
	config  *config.Config
	user    *user.Service
	order   *order.Service
	accrual *accrual.Service
	logger  *zap.SugaredLogger
}

func NewHandler(config *config.Config, user *user.Service, order *order.Service, accrual *accrual.Service, logger *zap.SugaredLogger) *Handler {
	return &Handler{config, user, order, accrual, logger}
}
