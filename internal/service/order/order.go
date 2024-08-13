package order

import (
	"context"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/storage/order"
	"github.com/mikesvis/gmart/pkg/luhn"
)

type Service struct {
	storage *order.Storage
}

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
var ErrConflict = errors.New(http.StatusText(http.StatusConflict))
var ErrExisted = errors.New(http.StatusText(http.StatusOK))
var ErrUnprocessableEntity = errors.New(http.StatusText(http.StatusUnprocessableEntity))
var ErrIternal = errors.New(http.StatusText(http.StatusInternalServerError))

func NewService(storage *order.Storage) *Service {
	return &Service{
		storage,
	}
}

func (s *Service) CreateOrder(ctx context.Context, orderID, userID uint64) error {
	if userID == 0 || orderID == 0 {
		return ErrBadRequest
	}

	if !luhn.IsValid(orderID) {
		return ErrUnprocessableEntity
	}

	err := s.storage.Create(ctx, orderID, userID, domain.StatusNew)
	// ошибка не по duplicate key
	if err != nil && !errors.Is(err, order.ErrConflict) {
		return err
	}

	// ошибки не было
	if err == nil {
		return nil
	}

	existingOrder, err := s.storage.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if existingOrder.UserID != userID {
		return ErrConflict
	}

	return ErrExisted
}
