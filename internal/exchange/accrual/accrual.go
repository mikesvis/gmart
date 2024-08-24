package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/domain"
	"go.uber.org/zap"
)

type Exchange struct {
	config *config.Config
	ch     chan *uint64
	logger *zap.SugaredLogger
}

type OrderAccrualResponse struct {
	OrderID int64         `json:"order,string"`
	Status  domain.Status `json:"status"`
	Accrual float64       `json:"accrual"`
}

func NewExchange(config *config.Config, logger *zap.SugaredLogger) *Exchange {
	ch := make(chan *uint64)

	// Я знаю это плохо и мне очень стыдно
	// Пока что организовал начисление в синхрооном режиме,
	// Сдаю сейчас на первичное ревью и боюсь что впринципе не успею, но буду стараться
	// Насчет написания тестов уже не уверен, мне бы асинхронность успеть :(
	// Прошу понять и простить, а пока приступаю к асинхронности, если будут замечания по тому что есть, тодже учту

	// go ProcessAccrual(ch, logger)

	return &Exchange{config, ch, logger}
}

func (e *Exchange) GetOrderAccrual(ctx context.Context, orderID uint64) (*OrderAccrualResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/orders/%d", e.config.AccrualSystemAddress, orderID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	e.logger.Infof("recieved accrual response %s status %d", body, resp.StatusCode)

	var orderAccrualResponse OrderAccrualResponse
	err = json.Unmarshal(body, &orderAccrualResponse)
	if err != nil {
		e.logger.Error(err)
		return nil, err
	}

	e.logger.Infof("decoded json %v", orderAccrualResponse)

	return &orderAccrualResponse, nil
}
