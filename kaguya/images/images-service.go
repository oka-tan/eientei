package images

import (
	"context"
	"fmt"
	"kaguya/api"
	"kaguya/config"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//QueuedImage is an image to be downloaded from 4chan and uploaded to S3.
type QueuedImage struct {
	Ext   string
	Tim   int64
	Board string
	Fsize int64
}

//Service is a singleton that obtains and uploads images to S3.
type Service struct {
	Queue      chan []QueuedImage
	host       string
	client     http.Client
	s3Client   *s3.Client
	napTime    time.Duration
	bucketName string
}

//NewService creates and returns a new images.Service
func NewService(
	imagesConfig config.ImagesConfig,
	s3Client *s3.Client,
) Service {
	requestTimeout, _ := time.ParseDuration(imagesConfig.RequestTimeout)

	client := http.Client{
		Timeout: requestTimeout,
	}

	napTime, _ := time.ParseDuration(imagesConfig.NapTime)

	return Service{
		Queue:      make(chan []QueuedImage, 200),
		host:       imagesConfig.Host,
		bucketName: imagesConfig.BucketName,
		client:     client,
		napTime:    napTime,
		s3Client:   s3Client,
	}
}

func (s *Service) Enqueue(boardName string, posts []api.Post) {
	if len(posts) > 0 {
		s.Queue <- toQueuedImages(boardName, posts)
	}
}

//Run is the main loop for images.Service instances.
func (s *Service) Run() {
	for {
		batch := <-s.Queue
		uploader := manager.NewUploader(s.s3Client)

		for _, image := range batch {
			file := fmt.Sprintf("%s/%d%s", image.Board, image.Tim, image.Ext)
			thumb := fmt.Sprintf("%s/%ds.jpg", image.Board, image.Tim)

			for _, f := range []string{file, thumb} {

				_, err := s.s3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
					Bucket: &s.bucketName,
					Key:    &f,
				})

				if err == nil {
					continue
				}

				time.Sleep(s.napTime)

				resp, err := s.client.Get(fmt.Sprintf("%s/%s", s.host, f))

				if err != nil {
					log.Printf("Error obtaining image:\nURL: %s\nError: %s\n", f, err)
					continue
				}

				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Printf("Error obtaining image:\nURL: %s\nError: status code is %d", f, resp.StatusCode)
					continue
				}

				contentType := resp.Header["Content-Type"][0]

				_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
					Bucket:      &s.bucketName,
					Key:         &f,
					Body:        resp.Body,
					ContentType: &contentType,
				})

				if err != nil {
					log.Printf("Error uploading file: %s", err)
				}
			}
		}
	}
}
