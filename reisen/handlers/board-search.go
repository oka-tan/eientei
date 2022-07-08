package handlers

import (
	"context"
	"log"
	"net/http"
	"reisen/config"
	"reisen/db"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type SearchRequest struct {
	Resto    int64  `query:"resto"`
	Com      string `query:"com"`
	Sub      string `query:"sub"`
	Name     string `query:"name"`
	Trip     string `query:"trip"`
	Filename string `query:"filename"`
	MD5      string `query:"md5"`
	Capcode  string `query:"capcode"`
	Image    string `query:"image"`
	Deleted  string `query:"deleted"`
	OP       string `query:"op"`
	Anon     string `query:"anon"`
	Keyset   int64  `query:"keyset"`
	RKeyset  int64  `query:"rkeyset"`
}

func (s SearchRequest) AllCapcodes() bool {
	return s.Capcode == ""
}

func (s SearchRequest) NoCapcode() bool {
	return s.Capcode == "null"
}

func (s SearchRequest) ModCapcode() bool {
	return s.Capcode == "mod"
}

func (s SearchRequest) ManagerCapcode() bool {
	return s.Capcode == "manager"
}

func (s SearchRequest) AdminCapcode() bool {
	return s.Capcode == "admin"
}

func (s SearchRequest) DeveloperCapcode() bool {
	return s.Capcode == "developer"
}

func (s SearchRequest) FounderCapcode() bool {
	return s.Capcode == "founder"
}

func (s SearchRequest) AllImages() bool {
	return s.Image == ""
}

func (s SearchRequest) HasImage() bool {
	return s.Image == "image"
}

func (s SearchRequest) NoImage() bool {
	return s.Image == "noImage"
}

func (s SearchRequest) Spoiler() bool {
	return s.Image == "spoiler"
}

func (s SearchRequest) NoSpoiler() bool {
	return s.Image == "noSpoiler"
}

func (s SearchRequest) AllDeleted() bool {
	return s.Deleted == ""
}

func (s SearchRequest) JustDeleted() bool {
	return s.Deleted == "true"
}

func (s SearchRequest) NoDeleted() bool {
	return s.Deleted == "false"
}

func (s SearchRequest) AllOPs() bool {
	return s.OP == ""
}

func (s SearchRequest) JustSticky() bool {
	return s.OP == "sticky"
}

func (s SearchRequest) NoSticky() bool {
	return s.OP == "noSticky"
}

func (s SearchRequest) JustOP() bool {
	return s.OP == "op"
}

func (s SearchRequest) JustReply() bool {
	return s.OP == "reply"
}

func (s SearchRequest) AllPosters() bool {
	return s.Anon == ""
}

func (s SearchRequest) Anonymous() bool {
	return s.Anon == "anon"
}

func (s SearchRequest) Namefag() bool {
	return s.Anon == "namefag"
}

func (s SearchRequest) Tripfag() bool {
	return s.Anon == "tripfag"
}

func (s SearchRequest) NameAndTripfag() bool {
	return s.Anon == "nameAndTripfag"
}

func BoardSearch(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		var result []*db.Post
		model := map[string]interface{}{
			"board":  board,
			"boards": conf.Boards,
			"conf":   conf.TemplateConfig,
		}

		req := new(SearchRequest)

		if err := c.Bind(req); err != nil {
			model["error"] = "Error binding request parameters."

			return c.Render(http.StatusOK, "board-search", model)
		}

		q := pg.NewSelect().
			Model(&result).
			Where("board = ?", board)

		if req.Resto > 0 {
			q.Where("resto = ?", req.Resto)
		}

		if req.Com != "" {
			q.Where("com_tsvector IS NOT NULL AND com_tsvector @@ plainto_tsquery('english', ?)", req.Com)
		}

		if req.Sub != "" {
			q.Where("OP AND sub_tsvector IS NOT NULL AND sub_tsvector @@ plainto_tsquery('english', ?)", req.Sub)
		}

		if req.Name != "" {
			q.Where("name_tsvector IS NOT NULL AND name_tsvector @@ plainto_tsquery('english', ?)", req.Name)
		}

		if req.Trip != "" {
			q.Where("trip = ?", req.Trip)
		}

		if req.Filename != "" {
			q.Where("tim IS NOT NULL AND filename_tsvector IS NOT NULL AND filename_tsvector @@ plainto_tsquery('english', ?)", req.Filename)
		}

		if req.MD5 != "" {
			q.Where("tim IS NOT NULL AND md5 = ?", req.MD5)
		}

		if req.Capcode != "" {
			if req.Capcode == "null" {
				q.Where("capcode IS NULL")
			} else {
				q.Where("capcode = ?", req.Capcode)
			}
		}

		if req.Image != "" {
			switch req.Image {
			case "image":
				{
					q.Where("tim IS NOT NULL")
				}
			case "noImage":
				{
					q.Where("tim IS NULL")
				}
			case "spoiler":
				{
					q.Where("tim IS NOT NULL AND spoiler")
				}
			case "noSpoiler":
				{
					q.Where("tim IS NOT NULL AND NOT spoiler")
				}
			}
		}

		if req.Deleted != "" {
			switch req.Deleted {
			case "true":
				{
					q.Where("deleted")
				}
			case "false":
				{
					q.Where("NOT deleted")
				}
			}
		}

		if req.OP != "" {
			switch req.OP {
			case "sticky":
				{
					q.Where("op AND sticky")
				}
			case "noSticky":
				{
					q.Where("op AND NOT sticky")
				}
			case "op":
				{
					q.Where("op")
				}
			case "reply":
				{
					q.Where("NOT op")
				}
			}
		}

		if req.Anon != "" {
			switch req.Anon {
			case "anon":
				{
					q.Where("name IS NULL AND trip IS NULL")
				}
			case "namefag":
				{
					q.Where("name IS NOT NULL")
				}
			case "tripfag":
				{
					q.Where("trip IS NOT NULL")
				}
			case "nameAndTripfag":
				{
					q.Where("name IS NOT NULL")
					q.Where("trip IS NOT NULL")
				}
			}
		}

		if req.Keyset > 0 {
			q.Where("no < ?", req.Keyset)
		} else if req.RKeyset > 0 {
			q.Where("no > ?", req.RKeyset)
		}

		q.Order("no DESC").Limit(24)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		if err := q.Scan(ctx); err != nil {
			log.Println(err)
			return c.Render(http.StatusOK, "board-search-error", model)
		}

		if len(result) == 24 {
			model["keyset"] = result[23].No
		}

		if req.Keyset != 0 || req.RKeyset != 0 {
			model["rkeyset"] = result[0].No
		}

		model["search"] = req
		model["result"] = result

		return c.Render(http.StatusOK, "board-search", model)
	}
}
