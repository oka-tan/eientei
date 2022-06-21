package handlers

import (
	"net/http"
	"reisen/config"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Index(db *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		model := map[string]interface{}{
			"boards": conf.Boards,
			"conf":   conf.TemplateConfig,
		}

		return c.Render(http.StatusOK, "index", model)
	}
}
