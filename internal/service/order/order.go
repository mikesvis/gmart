package order

import (
	"context"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/storage/order"
	"github.com/mikesvis/gmart/pkg/luhn"
	"go.uber.org/zap"
)

type Service struct {
	storage *order.Storage
	logger  *zap.SugaredLogger
}

var ErrNoContent = errors.New(http.StatusText(http.StatusNoContent))
var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
var ErrConflict = errors.New(http.StatusText(http.StatusConflict))
var ErrExisted = errors.New(http.StatusText(http.StatusOK))
var ErrUnprocessableEntity = errors.New(http.StatusText(http.StatusUnprocessableEntity))
var ErrIternal = errors.New(http.StatusText(http.StatusInternalServerError))

func NewService(storage *order.Storage, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage,
		logger,
	}
}

func (s *Service) CreateOrder(ctx context.Context, orderID, userID uint64) error {
	if userID == 0 || orderID == 0 {
		return ErrBadRequest
	}

	if !luhn.IsValid(orderID) {
		s.logger.Warnf("error while validating order number for %d", orderID)
		return ErrUnprocessableEntity
	}

	s.logger.Infof("creating order %d for user %d", orderID, userID)
	err := s.storage.Create(ctx, orderID, userID, domain.StatusNew)
	// ошибка не по duplicate key
	if err != nil && !errors.Is(err, order.ErrConflict) {
		s.logger.Errorf("error while creating order %d for user %d", orderID, userID, err)
		return err
	}

	// ошибки не было
	if err == nil {
		s.logger.Infof("success in creating order %d for user %d", orderID, userID)
		return nil
	}

	s.logger.Infof("searching existing order %d for user %d", orderID, userID)
	existingOrder, err := s.storage.FindByID(ctx, orderID)
	if err != nil {
		s.logger.Errorf("error in searching order %d for user %d", orderID, userID, err)
		return err
	}

	if existingOrder.UserID != userID {
		s.logger.Warnf("order %d is now owned by user %d", orderID, userID)
		return ErrConflict
	}

	s.logger.Infof("order %d is now already existed and owned by user %d", orderID, userID)
	return ErrExisted
}

func (s *Service) GetOrdersByUser(ctx context.Context, userID uint64) ([]domain.Order, error) {
	if userID == 0 {
		return nil, ErrBadRequest
	}

	s.logger.Infof("searching orders for user %d", userID)
	orders, err := s.storage.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("error in searching orders for user %d", userID, err)
		return nil, err
	}

	if len(orders) == 0 {
		s.logger.Infof("no orders for user %d", userID)
		return nil, ErrNoContent
	}

	s.logger.Infof("found %d orders for user %d", len(orders), userID)
	return orders, nil
}
