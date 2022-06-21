package handlers

import (
	"context"
	"net/http"
	"reisen/config"
	"reisen/db"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func BoardThread(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		no := c.Param("no")

		var thread []*db.Post

		err := pg.NewSelect().
			Model(&thread).
			Where("board = ? AND resto = ?", board, no).
			Order("no ASC").
			Scan(context.Background())

		if err != nil {
			model := map[string]interface{}{
				"board":  board,
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"no":     no,
			}

			return c.Render(http.StatusOK, "board-thread-error", model)
		}

		if len(thread) == 0 {
			model := map[string]interface{}{
				"board":  board,
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"no":     no,
			}

			return c.Render(http.StatusOK, "board-thread-not-found", model)
		}

		model := map[string]interface{}{
			"board":  board,
			"boards": conf.Boards,
			"conf":   conf.TemplateConfig,
			"op":     thread[0],
			"thread": thread[1:],
		}

		return c.Render(http.StatusOK, "board-thread", model)
	}
}
