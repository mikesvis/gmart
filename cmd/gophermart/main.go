package main

import (
	"context"

	"github.com/mikesvis/gmart/internal/app/gophermart"
	"github.com/mikesvis/gmart/internal/config"
)

func main() {
	ctx := context.Background()
	config := config.NewConfig()
	gophermart, err := gophermart.NewGophermart(config)
	if err != nil {
		panic(err)
	}

	gophermart.Run(ctx)
}
