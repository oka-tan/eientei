package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type BoardRequest struct {
	Board  string `param:"board"`
	Keyset int64  `query:"keyset"`
	Order  string `query:"order"`
}

func BoardHandler(db *bun.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var boardRequest BoardRequest
		if err := c.Bind(&boardRequest); err != nil {
			return err
		}

		var threads []Post
		q := db.NewSelect().
			Model(&threads).
			Where("board = ?", boardRequest.Board).
			Where("op")

		if boardRequest.Keyset > 0 {
			if boardRequest.Order == "asc" {
				q.Where("no > ?", boardRequest.Keyset)
			} else {
				q.Where("no < ?", boardRequest.Keyset)
			}
		}

		if boardRequest.Order == "asc" {
			q.Order("no ASC")
		} else {
			q.Order("no DESC")
		}

		q.Limit(200)

		if err := q.Scan(context.Background()); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if len(threads) == 0 {
			return c.NoContent(http.StatusNoContent)
		}

		return c.JSON(http.StatusOK, threads)
	}
}
