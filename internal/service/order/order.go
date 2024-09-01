package order

import (
	"context"
	"errors"

	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/storage/order"
	"github.com/mikesvis/gmart/pkg/luhn"
	"go.uber.org/zap"
)

type Service struct {
	storage *order.Storage
	logger  *zap.SugaredLogger
}

var ErrNoOrdersInSet = errors.New(`orders list is empty`)
var ErrWrongArgument = errors.New(`argument can not equal to 0`)
var ErrOwnerIsIncorrect = errors.New(`order is owned by another user`)
var ErrOwnedByAnother = errors.New(`order exists and owned by another user`)
var ErrOrderNumber = errors.New(`order number is invalid`)

func NewService(storage *order.Storage, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage,
		logger,
	}
}

func (s *Service) CreateOrder(ctx context.Context, orderID, userID uint64) error {
	if userID == 0 || orderID == 0 {
		return ErrWrongArgument
	}

	if !luhn.IsValid(orderID) {
		s.logger.Warnf("error while validating order number for %d", orderID)
		return ErrOrderNumber
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
		return ErrOwnerIsIncorrect
	}

	s.logger.Infof("order %d is already existed and owned by user %d", orderID, userID)
	return ErrOwnedByAnother
}

func (s *Service) GetOrdersByUser(ctx context.Context, userID uint64) ([]domain.Order, error) {
	if userID == 0 {
		return nil, ErrWrongArgument
	}

	s.logger.Infof("searching orders for user %d", userID)
	orders, err := s.storage.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("error in searching orders for user %d", userID, err)
		return nil, err
	}

	if len(orders) == 0 {
		s.logger.Infof("no orders for user %d", userID)
		return nil, ErrNoOrdersInSet
	}

	s.logger.Infof("found %d orders for user %d", len(orders), userID)
	return orders, nil
}
