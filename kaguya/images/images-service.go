package images

import (
	"context"
	"fmt"
	"kaguya/api"
	"kaguya/config"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type queuedImage struct {
	Ext   string
	Tim   int64
	Board string
}

//Service is a singleton that obtains and uploads images to S3.
type Service struct {
	Queue    chan queuedImage
	host     string
	client   http.Client
	s3Client *s3.Client
	napTime  time.Duration
	bucket   string
}

//NewService creates and returns a new images.Service
func NewService(
	imagesConfig config.ImagesConfig,
) Service {
	requestTimeout, _ := time.ParseDuration(imagesConfig.RequestTimeout)

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if imagesConfig.S3Host != "" {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               imagesConfig.S3Host,
				SigningRegion:     imagesConfig.Region,
				HostnameImmutable: true,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	s3Config, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithEndpointResolver(customResolver),
		awsconfig.WithRegion(imagesConfig.Region),
	)

	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(s3Config)

	client := http.Client{
		Timeout: requestTimeout,
	}

	napTime, _ := time.ParseDuration(imagesConfig.NapTime)

	return Service{
		Queue:    make(chan queuedImage, 1000000),
		host:     imagesConfig.Host,
		bucket:   imagesConfig.Bucket,
		client:   client,
		napTime:  napTime,
		s3Client: s3Client,
	}
}

func (s *Service) Enqueue(boardName string, posts []api.Post) {
	for _, p := range posts {
		if p.Tim != nil && p.Ext != nil {
			s.Queue <- queuedImage{
				Tim:   *(p.Tim),
				Ext:   *(p.Ext),
				Board: boardName,
			}
		}
	}
}

func (s *Service) EnqueueMap(boardName string, posts map[int64]api.Post) {
	for _, p := range posts {
		if p.Tim != nil && p.Ext != nil {
			s.Queue <- queuedImage{
				Tim:   *(p.Tim),
				Ext:   *(p.Ext),
				Board: boardName,
			}
		}
	}
}

//Run is the main loop for images.Service instances.
func (s *Service) Run() {
	uploader := manager.NewUploader(s.s3Client, func(u *manager.Uploader) {
		//Theoretically, the file size cap is 6mb because of /wsg/
		//We add in 2mb to be safe and then 2 extra mb on the buffer to be safer
		u.PartSize = 8 * 1024 * 1024
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(10 * 1024 * 1024)
	})

	for image := range s.Queue {
		file := fmt.Sprintf("%s/%d%s", image.Board, image.Tim, image.Ext)

		time.Sleep(s.napTime)

		resp, err := s.client.Get(fmt.Sprintf("%s/%s", s.host, file))

		if err != nil {
			log.Printf("Error obtaining image:\nURL: %s\nError: %s\n", file, err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("Error obtaining image:\nURL: %s\nError: status code is %d", file, resp.StatusCode)
			continue
		}

		contentType := resp.Header["Content-Type"][0]

		_, err = uploader.Upload(context.Background(), &s3.PutObjectInput{
			Bucket:      &s.bucket,
			Key:         &file,
			Body:        resp.Body,
			ContentType: &contentType,
		})

		if err != nil {
			log.Printf("Error uploading file: %s", err)
		}
	}
}
