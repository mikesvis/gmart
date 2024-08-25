package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mikesvis/gmart/internal/config"
	driversDb "github.com/mikesvis/gmart/internal/drivers/postgres"
	driversQueue "github.com/mikesvis/gmart/internal/drivers/queue"
	accrualExchange "github.com/mikesvis/gmart/internal/exchange/accrual"
	"github.com/mikesvis/gmart/internal/handler"
	"github.com/mikesvis/gmart/internal/logger"
	server "github.com/mikesvis/gmart/internal/router"
	accrualService "github.com/mikesvis/gmart/internal/service/accrual"
	orderService "github.com/mikesvis/gmart/internal/service/order"
	userService "github.com/mikesvis/gmart/internal/service/user"
	accrualStorage "github.com/mikesvis/gmart/internal/storage/accrual"
	orderStorage "github.com/mikesvis/gmart/internal/storage/order"
	userStorage "github.com/mikesvis/gmart/internal/storage/user"
	"go.uber.org/zap"
)

type App struct {
	config *config.Config
	logger *zap.SugaredLogger
	router *chi.Mux
}

func New(ctx *context.Context) *App {
	config := config.NewConfig()
	logger := logger.NewLogger()
	db, _ := driversDb.NewPostgres(config)
	queue := driversQueue.NewQueue()

	userStorage := userStorage.NewStorage(db, logger)
	userService := userService.NewService(userStorage, logger)

	orderStorage := orderStorage.NewStorage(db, logger)
	orderService := orderService.NewService(orderStorage, logger)

	accrualStorage := accrualStorage.NewStorage(db, logger)
	accrualExchange := accrualExchange.NewExchange(config, logger)
	accrualService := accrualService.NewService(accrualStorage, accrualExchange, queue, logger)
	go accrualService.RunQueue(*ctx)

	handler := handler.NewHandler(config, userService, orderService, accrualService, logger)
	router := server.NewRouter(handler)
	return &App{
		config,
		logger,
		router,
	}
}

func (app *App) Run() {
	app.logger.Infow("Config initialized", "config", app.config)
	if err := http.ListenAndServe(string(app.config.RunAddress), app.router); err != nil {
		app.logger.Fatalw(err.Error(), "event", "start server")
	}
}
