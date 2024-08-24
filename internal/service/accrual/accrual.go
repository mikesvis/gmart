package accrual

import (
	"context"
	"errors"
	"net/http"

	"github.com/mikesvis/gmart/internal/domain"
	accrualExchange "github.com/mikesvis/gmart/internal/exchange/accrual"
	accrualStorage "github.com/mikesvis/gmart/internal/storage/accrual"
	"go.uber.org/zap"
)

type Service struct {
	storage  *accrualStorage.Storage
	exchange *accrualExchange.Exchange
	logger   *zap.SugaredLogger
}

var ErrBadRequest = errors.New(http.StatusText(http.StatusBadRequest))
var ErrNoContent = errors.New(http.StatusText(http.StatusNoContent))

func NewService(storage *accrualStorage.Storage, exchange *accrualExchange.Exchange, logger *zap.SugaredLogger) *Service {
	return &Service{
		storage,
		exchange,
		logger,
	}
}

func (s *Service) GetUserBalance(ctx context.Context, userID uint64) (*domain.UserBalance, error) {
	if userID == 0 {
		return nil, ErrBadRequest
	}

	s.logger.Infof("getting user balance for user %d", userID)
	balance, err := s.storage.GetBalanceByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("error while getting user balance for user %d, %v", userID, err)
		return nil, err
	}
	s.logger.Infof("user balance for user %d is %v", userID, balance)

	return balance, nil
}

func (s *Service) WithdrawToOrderID(ctx context.Context, orderID uint64, sum int64) error {
	if orderID == 0 {
		return ErrBadRequest
	}

	s.logger.Infof("creating withdraw for order %d", orderID)
	err := s.storage.CreateWithdrawn(ctx, orderID, sum, domain.StatusProcessed)
	if err != nil {
		s.logger.Errorf("error while creating withdraw for order %d, %v", orderID, err)
	}

	return err
}

func (s *Service) GetUserWithdrawals(ctx context.Context, userID uint64) ([]domain.UserWithdrawals, error) {
	if userID == 0 {
		return nil, ErrBadRequest
	}

	s.logger.Infof("getting user withdrawals for user %d", userID)
	withdrawals, err := s.storage.GetUserWithdrawals(ctx, userID)
	if err != nil {
		s.logger.Errorf("error while getting withdrawals for user %d, %v", userID, err)
		return nil, err
	}

	if len(withdrawals) == 0 {
		s.logger.Infof("empty withdrawals for user %d", userID)
		return nil, ErrNoContent
	}

	return withdrawals, nil
}

func (s *Service) ProcessOrderAccrual(ctx context.Context, orderID uint64) error {
	s.logger.Infof("processing accrual for order %d", orderID)
	orderAccrual, err := s.exchange.GetOrderAccrual(ctx, orderID)
	if err != nil {
		s.logger.Errorf("error while processing accrual for order %d, %v", orderID, err)
		return err
	}

	if orderAccrual == nil {
		s.logger.Infof("empty accural for order %d", orderID)
		return ErrNoContent
	}

	s.logger.Infof("creating accrual for order %d, %v", orderID, orderAccrual)
	err = s.storage.CreateAccrual(ctx, orderID, orderAccrual.Status, uint64(orderAccrual.Accrual*100))
	if err != nil {
		s.logger.Infof("error while creating accrual for order %d, %v, %v", orderID, orderAccrual, err)
		return err
	}

	orderStatus := domain.StatusProcessing
	if orderAccrual.Status != domain.StatusRegistered {
		orderStatus = orderAccrual.Status
	}

	s.logger.Infof("updating order %d to status %s", orderID, orderStatus)
	err = s.storage.UpdateOrderStatus(ctx, orderID, orderStatus)
	if err != nil {
		s.logger.Errorf("error while updating order %d to status %s", orderID, orderStatus, err)
	}

	return err
}
