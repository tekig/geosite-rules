package delivery

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tekig/geosite-rules/internal/repository"
	"golang.org/x/sync/errgroup"
)

type HTTP struct {
	e      *echo.Echo
	repo   repository.Repository
	listen string
}

func NewHTTP(listen string, repo repository.Repository) (*HTTP, error) {
	h := &HTTP{
		e:      echo.New(),
		repo:   repo,
		listen: listen,
	}

	h.e.GET("/rule/geosite/*", h.handler)

	h.e.Use(
		middleware.Logger(),
	)

	return h, nil
}

func (h *HTTP) handler(c echo.Context) error {
	rule := c.Param("*")

	rule = strings.TrimSuffix(rule, ".list")

	rules, err := h.repo.Rules(c.Request().Context(), rule)
	if err != nil {
		return fmt.Errorf("rules: %w", err)
	}

	return c.Blob(http.StatusOK, echo.MIMETextPlainCharsetUTF8, []byte(rules))
}

func (h *HTTP) Run(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return h.e.Start(h.listen)
	})

	wg.Go(func() error {
		<-ctx.Done()

		h.e.Close()

		return nil
	})

	return wg.Wait()
}
