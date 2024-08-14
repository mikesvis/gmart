package accural

import (
	"context"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/domain"
	"github.com/mikesvis/gmart/internal/storage/accural"
)

type Service struct {
	storage *accural.Storage
}

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))

func NewService(storage *accural.Storage) *Service {
	return &Service{
		storage,
	}
}

func (s *Service) GetUserBalance(ctx context.Context, userID uint64) (*domain.UserBalance, error) {
	if userID == 0 {
		return nil, ErrBadRequest
	}

	balance, err := s.storage.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
