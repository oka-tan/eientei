package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type PostsRequest struct {
	Board  string `param:"board"`
	Keyset int64  `query:"keyset"`
	Order  string `query:"order"`
}

func PostsHandler(db *bun.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		var postsRequest PostsRequest
		if err := c.Bind(&postsRequest); err != nil {
			return err
		}

		var posts []Post
		q := db.NewSelect().
			Model(&posts).
			Where("board = ?", postsRequest.Board)

		if postsRequest.Keyset > 0 {
			if postsRequest.Order == "asc" {
				q.Where("no > ?", postsRequest.Keyset)
			} else {
				q.Where("no < ?", postsRequest.Keyset)
			}
		}

		if postsRequest.Order == "asc" {
			q.Order("no ASC")
		} else {
			q.Order("no DESC")
		}

		q.Limit(200)

		if err := q.Scan(context.Background()); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if len(posts) == 0 {
			return c.NoContent(http.StatusNoContent)
		}

		return c.JSON(http.StatusOK, posts)
	}
}
