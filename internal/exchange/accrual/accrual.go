package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mikesvis/gmart/internal/config"
	"github.com/mikesvis/gmart/internal/domain"
	"go.uber.org/zap"
)

type AccrualExchange struct {
	config *config.Config
	logger *zap.SugaredLogger
}

type orderAccrualResponse struct {
	OrderID int64         `json:"order,string"`
	Status  domain.Status `json:"status"`
	Accrual float64       `json:"accrual"`
}

var ErrRequeue = errors.New("requeue to accrual")

func NewAccrualExchange(config *config.Config, logger *zap.SugaredLogger) *AccrualExchange {
	return &AccrualExchange{config, logger}
}

func (e *AccrualExchange) GetOrderAccrual(ctx context.Context, orderID uint64) (*orderAccrualResponse, error) {
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

	if resp.StatusCode == http.StatusTooManyRequests {
		defaultSleepDuration := int64(60)
		resposeSleepDuration, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			resposeSleepDuration = defaultSleepDuration
		}
		e.logger.Infof("recieved accrual response status %d sleeping for %d seconds", resp.StatusCode, resposeSleepDuration)
		time.Sleep(time.Second * time.Duration(resposeSleepDuration))

		return nil, ErrRequeue
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	e.logger.Infof("recieved accrual response %s status %d", body, resp.StatusCode)

	var orderAccrualResponse orderAccrualResponse
	err = json.Unmarshal(body, &orderAccrualResponse)
	if err != nil {
		e.logger.Error(err)
		return nil, err
	}

	e.logger.Infof("decoded json %v", orderAccrualResponse)

	return &orderAccrualResponse, nil
}
