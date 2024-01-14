package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var _ Repository = (*Plain)(nil)

type Plain struct {
	path string
}

func NewPlain(path string) (*Plain, error) {
	return &Plain{
		path: path,
	}, nil
}

func (r *Plain) Rules(ctx context.Context, name string) (string, error) {
	path := filepath.Join(r.path, filepath.Clean(name))

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrRuleNotFound
	} else if err != nil {
		return "", fmt.Errorf("read file")
	}

	return string(data), nil
}
