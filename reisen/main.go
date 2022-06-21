package main

import (
	"database/sql"
	"fmt"
	"reisen/config"
	"reisen/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	conf, err := config.LoadConfig()

	if err != nil {
		panic(fmt.Sprintf("Error loading configuration: %s", err))
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	db := bun.NewDB(sqldb, pgdialect.New())

	e := echo.New()
	e.Renderer = NewTemplater()
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		XFrameOptions:         "",
		ContentSecurityPolicy: conf.CspConfig,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  log.ERROR,
	}))
	e.Use(middleware.Gzip())

	e.GET("/", handlers.Index(db, conf))
	e.GET("/contact", handlers.Contact(db, conf))
	e.GET("/:board", handlers.Board(db, conf))
	e.GET("/:board/search", handlers.BoardSearch(db, conf))
	e.GET("/:board/thread/:no", handlers.BoardThread(db, conf))

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
