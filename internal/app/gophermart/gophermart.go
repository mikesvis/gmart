package gophermart

import (
	"context"
	"net/http"

	"github.com/mikesvis/gmart/internal/config"
	accrualExchange "github.com/mikesvis/gmart/internal/exchange/accrual"
	"github.com/mikesvis/gmart/internal/handler"
	"github.com/mikesvis/gmart/internal/logger"
	"github.com/mikesvis/gmart/internal/postgres"
	"github.com/mikesvis/gmart/internal/queue"
	"github.com/mikesvis/gmart/internal/server"

	accrualService "github.com/mikesvis/gmart/internal/service/accrual"
	orderService "github.com/mikesvis/gmart/internal/service/order"
	userService "github.com/mikesvis/gmart/internal/service/user"
	accrualStorage "github.com/mikesvis/gmart/internal/storage/accrual"
	orderStorage "github.com/mikesvis/gmart/internal/storage/order"
	userStorage "github.com/mikesvis/gmart/internal/storage/user"
	"go.uber.org/zap"
)

type Gophermart struct {
	config         *config.Config
	logger         *zap.SugaredLogger
	server         *http.Server
	accrualService *accrualService.Service
}

func NewGophermart(config *config.Config) (*Gophermart, error) {

	logger := logger.NewLogger()

	db, err := postgres.NewPostgres(config)
	if err != nil {
		return nil, err
	}

	queue := queue.NewQueue()

	userStorage := userStorage.NewStorage(db, logger)
	userService := userService.NewService(userStorage, logger)

	orderStorage := orderStorage.NewStorage(db, logger)
	orderService := orderService.NewService(orderStorage, logger)

	accrualStorage := accrualStorage.NewStorage(db, logger)
	accrualExchange := accrualExchange.NewAccrualExchange(config, logger)
	accrualService := accrualService.NewService(accrualStorage, accrualExchange, queue, logger)

	handler := handler.NewHandler(config, userService, orderService, accrualService, logger)
	router := server.NewRouter(handler)
	server := &http.Server{Addr: config.RunAddress, Handler: router}
	return &Gophermart{
		config,
		logger,
		server,
		accrualService,
	}, nil
}

func (g *Gophermart) Run(ctx context.Context) {
	g.logger.Infow("Run with config", "config", g.config)
	go g.accrualService.RunQueue(ctx)
	if err := g.server.ListenAndServe(); err != nil {
		g.logger.Fatalw(err.Error(), "event", "start server")
	}
}
