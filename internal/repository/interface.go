package repository

import (
	"context"
)

type Repository interface {
	Rules(ctx context.Context, name string) (string, error)
}
