package app

import (
	"context"
	"fmt"
	"os"

	"github.com/tekig/geosite-rules/internal/delivery"
	"github.com/tekig/geosite-rules/internal/repository"
)

type App struct {
	d *delivery.HTTP
}

func New() (*App, error) {
	r, err := repository.NewDLC(env("DLC_PATH", "dlc.dat"))
	if err != nil {
		return nil, fmt.Errorf("dlc: %w", err)
	}

	d, err := delivery.NewHTTP(env("LISTEN", ":80"), r)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}

	return &App{
		d: d,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.d.Run(ctx)
}

func env(name, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}

	return v
}
