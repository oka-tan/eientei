package main

import (
	"context"
	"database/sql"
	"log"
	"moon/config"
	"moon/db"
	"moon/lnx"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	time.Sleep(5 * time.Second)

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	d := bun.NewDB(sqldb, pgdialect.New())

	dbPosts := make([]db.Post, 0, conf.BatchSize)

	lnxService := lnx.NewService(conf.LnxConfig.Host, conf.LnxConfig.Port)

	//We undo all previous partial changes in case it crashed
	for _, board := range conf.Boards {
		if err := lnxService.Rollback(board); err != nil {
			panic(err)
		}
	}

	for {
		for _, board := range conf.Boards {
			log.Printf("Indexing board %s\n", board)

			indexTracker := db.IndexTracker{
				Board:        board,
				LastModified: time.UnixMicro(0), //Convenient old date.
				No:           0,
			}

			//We insert it beforehand in case it doesn't exist
			_, err := d.NewInsert().
				Model(&indexTracker).
				On("CONFLICT (board) DO NOTHING").
				Returning("NULL").
				Exec(context.Background())

			//We select it from the database in case it does exist
			err = d.NewSelect().
				Model(&indexTracker).
				Where("board = ?", board).
				Scan(context.Background())

			if err != nil {
				panic(err)
			}

			for {

				dbPosts = dbPosts[0:0]
				err := d.NewSelect().
					Model(&dbPosts).
					Where("board = ?", board).
					Where("(last_modified, no) > (?, ?)", indexTracker.LastModified, indexTracker.No).
					Order("last_modified ASC", "no ASC").
					Limit(int(conf.BatchSize)).
					Scan(context.Background())

				if err != nil {
					panic(err)
				}

				if len(dbPosts) == 0 {
					break
				} else {
					lastPost := dbPosts[len(dbPosts)-1]
					indexTracker.LastModified = lastPost.LastModified
					indexTracker.No = lastPost.No
				}

				if err = lnxService.Upsert(dbPosts, board); err != nil {
					panic(err)
				}
			}

			if err := lnxService.Commit(board); err != nil {
				panic(err)
			}

			_, err = d.NewUpdate().
				Model(&indexTracker).
				WherePK().
				Exec(context.Background())

			if err != nil {
				panic(err)
			}
		}

		log.Println("Napping")
		time.Sleep(20 * time.Minute)
	}
}
