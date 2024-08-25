package main

import (
	"context"

	"github.com/mikesvis/gmart/internal/app"
)

func main() {
	ctx := context.Background()
	app := app.New(&ctx)

	app.Run()
}
