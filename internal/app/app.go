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
	repoDLC, err := repository.NewDLC(env("DLC_PATH", "dlc.dat"))
	if err != nil {
		return nil, fmt.Errorf("dlc: %w", err)
	}

	repoPlain, err := repository.NewPlain(env("PLAIN_PATH", "./plain"))
	if err != nil {
		return nil, fmt.Errorf("plain: %w", err)
	}

	d, err := delivery.NewHTTP(
		env("LISTEN", ":80"),
		repository.NewMulti(
			repoDLC,
			repoPlain,
		),
	)
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
