package logger

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger.Sugar()
}
