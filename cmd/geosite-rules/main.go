package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/tekig/geosite-rules/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := a.Run(ctx); err != nil {
		println(err.Error())
		os.Exit(2)
	}
}
