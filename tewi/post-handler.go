package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type PostRequest struct {
	Board string `param:"board"`
	No    int64  `param:"no"`
}

func PostHandler(db *bun.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var postRequest PostRequest
		if err := (&echo.DefaultBinder{}).BindPathParams(c, &postRequest); err != nil {
			return err
		}

		var post Post
		err := db.NewSelect().
			Model(&post).
			Where("no = ?", postRequest.No).
			Where("board = ?", postRequest.Board).
			Scan(context.Background())

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if post.No == 0 {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, post)
	}
}
