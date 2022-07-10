package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type ThreadRequest struct {
	Board string `param:"board"`
	Resto int64  `param:"resto"`
}

func ThreadHandler(db *bun.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var threadRequest ThreadRequest

		if err := (&echo.DefaultBinder{}).BindPathParams(c, &threadRequest); err != nil {
			return err
		}

		var posts []Post
		err := db.NewSelect().
			Model(&posts).
			Where("resto = ?", threadRequest.Resto).
			Where("board = ?", threadRequest.Board).
			Order("no ASC").
			Scan(context.Background())

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if len(posts) == 0 {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, posts)
	}
}
