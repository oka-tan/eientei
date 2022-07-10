package main

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	conf, err := LoadConfig()

	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %s", err))
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	db := bun.NewDB(sqldb, pgdialect.New())

	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))
	e.Use(middleware.Gzip())

	e.GET("/:board", BoardHandler(db))
	e.GET("/:board/thread/:resto", ThreadHandler(db))
	e.GET("/:board/post/:no", PostHandler(db))
	e.GET("/:board/posts", PostsHandler(db))

	if conf.Production {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(conf.Hosts...)
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Pre(middleware.HTTPSNonWWWRedirect())

		go func() {
			e2 := echo.New()
			e2.Pre(middleware.HTTPSNonWWWRedirect())
			e2.Pre(middleware.HTTPSRedirect())

			e2.Logger.Fatal(e2.Start(":80"))
		}()

		e.Logger.Fatal(e.StartAutoTLS(":443"))
	} else {
		e.Logger.Fatal(e.Start(":1323"))
	}
}
