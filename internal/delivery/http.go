package delivery

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tekig/geosite-rules/internal/repository"
	"golang.org/x/sync/errgroup"
)

type HTTP struct {
	e      *echo.Echo
	repo   *repository.DLC
	listen string
}

func NewHTTP(listen string, repo *repository.DLC) (*HTTP, error) {
	h := &HTTP{
		e:      echo.New(),
		repo:   repo,
		listen: listen,
	}

	h.e.GET("/rule/geosite/*", h.handler)

	return h, nil
}

func (h *HTTP) handler(c echo.Context) error {
	rule := c.Param("*")

	rule = strings.TrimSuffix(rule, ".list")

	rules, err := h.repo.Rules(c.Request().Context(), rule)
	if err != nil {
		return fmt.Errorf("rules: %w", err)
	}

	list, err := repository.PlainText(rules)
	if err != nil {
		return fmt.Errorf("plain text: %w", err)
	}

	return c.Blob(http.StatusOK, echo.MIMETextPlainCharsetUTF8, []byte(list))
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
