package main

import (
	"context"
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

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	runtime.GOMAXPROCS(1)
	log.Println("Starting Kaguya")

	conf, err := config.LoadConfig()

	if err != nil {
		log.Panicf("Error loading configuration: %s", err)
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(conf.ImagesConfig.AwsRegion))

	if err != nil {
		log.Panicf("Error loading AWS configuration: %s", err)
	}

	s3Client := s3.NewFromConfig(awsConfig)

	imagesService := images.NewService(conf.ImagesConfig, s3Client)
	thumbnailsService := thumbnails.NewService(conf.ThumbnailsConfig, s3Client)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqldb, pgdialect.New())

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
	}

	go func(imagesService images.Service) {
		imagesService.Run()
	}(imagesService)

	go func(thumbnailsService thumbnails.Service) {
		thumbnailsService.Run()
	}(thumbnailsService)

	//Kaguya hoards substantially more memory than it actually uses, for some reason.
	//It's a GC tuning problem and not something I can actually solve by just rewriting the code.
	//Except golang doesn't really allow for any GC tuning other than this.
	go func() {
		for {
			runtime.GC()
			time.Sleep(time.Minute)
		}
	}()

	forever := make(chan bool)
	<-forever
}
