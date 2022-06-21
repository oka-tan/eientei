package handlers

import (
	"context"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func Board(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {

		board := c.Param("board")

		var threads []*db.Post

		q := pg.NewSelect().Model(&threads)
		keyset, kerr := strconv.ParseInt(c.QueryParam("keyset"), 10, 64)
		rkeyset, rkerr := strconv.ParseInt(c.QueryParam("rkeyset"), 10, 64)

		if kerr == nil {
			q.Where("OP AND board = ? AND no < ?", board, keyset)
		} else if rkerr == nil {
			q.Where("OP AND board = ? AND no > ?", board, rkeyset)
		} else {
			q.Where("OP AND board = ?", board)
		}

		err := q.
			Order("no DESC").
			Limit(24).
			Scan(context.Background())

		if err != nil {
			return c.Render(http.StatusOK, "board-error", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
			})
		}

		if len(threads) == 0 {
			return c.Render(http.StatusOK, "board-empty", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
			})
		}

		if keyset != 0 || rkeyset != 0 {
			rkeyset = threads[0].No
		} else {
			rkeyset = 0
		}

		if len(threads) == 24 {
			keyset = threads[len(threads)-1].No
		} else {
			keyset = 0
		}

		model := map[string]interface{}{
			"boards":  conf.Boards,
			"conf":    conf.TemplateConfig,
			"board":   board,
			"threads": threads,
			"keyset":  keyset,
			"rkeyset": rkeyset,
		}

		return c.Render(http.StatusOK, "board", model)
	}
}
