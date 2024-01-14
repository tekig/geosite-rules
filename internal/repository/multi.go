package repository

import (
	"context"
	"errors"
	"fmt"
)

type Multi struct {
	r []Repository
}

func NewMulti(r ...Repository) *Multi {
	return &Multi{
		r: r,
	}
}

func (r *Multi) Rules(ctx context.Context, name string) (string, error) {
	for _, r := range r.r {
		rule, err := r.Rules(ctx, name)
		if errors.Is(err, ErrRuleNotFound) {
			continue
		} else if err != nil {
			return "", fmt.Errorf("%T: %w", r, err)
		}

		return rule, nil
	}

	return "", ErrRuleNotFound
}
