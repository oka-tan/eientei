package main

import (
	"database/sql"
	"kaguya/config"
	"kaguya/images"
	"kaguya/manager"
	"kaguya/thumbnails"
	"log"
	"runtime"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	runtime.GOMAXPROCS(1)
	log.Println("Starting Kaguya")

	conf, err := config.LoadConfig()

	if err != nil {
		log.Panicf("Error loading configuration: %s", err)
	}

	var imagesService *images.Service

	if conf.StartImagesService() {
		imagesServiceP := images.NewService(conf.ImagesConfig)
		imagesService = &imagesServiceP
	}

	var thumbnailsService *thumbnails.Service
	if conf.StartThumbnailsService() {
		thumbnailsServiceP := thumbnails.NewService(conf.ThumbnailsConfig)
		thumbnailsService = &thumbnailsServiceP
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqldb, pgdialect.New())
	initialNap, _ := time.ParseDuration(conf.InitialNap)

	if imagesService != nil {
		go func(imagesService *images.Service) {
			imagesService.Run()
		}(imagesService)
	}

	if thumbnailsService != nil {
		go func(thumbnailsService *thumbnails.Service) {
			thumbnailsService.Run()
		}(thumbnailsService)
	}

	for _, boardConfig := range conf.Boards {
		go func(apiConfig config.APIConfig, boardConfig config.BoardConfig, pg *bun.DB) {

			log.Printf("Creating board manager for %+v", boardConfig)

			boardManager, err := manager.NewBoardManager(apiConfig, boardConfig, pg, imagesService, thumbnailsService)

			if err != nil {
				log.Panicf("Error initiating board manager:\nBoard: %s\nError: %s", boardConfig.Name, err)
			}

			log.Printf("Initing board manager for board %s", boardConfig.Name)

			if err := boardManager.Init(conf.SkipArchive); err != nil {
				panic(err)
			}

			log.Printf("Running boardManager for board %s", boardConfig.Name)

			if err := boardManager.Run(); err != nil {
				panic(err)
			}

		}(conf.APIConfig, boardConfig, pg)

		time.Sleep(initialNap)
	}

	forever := make(chan bool)
	<-forever
}
