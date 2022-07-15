package lnx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"moon/db"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	host   string
	client http.Client
}

func NewService(host string, port uint64) Service {
	return Service{
		host: host + ":" + strconv.FormatUint(port, 10) + "/indexes",
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) Upsert(posts []db.Post, board string) error {
	deletables := make([]db.Post, 0, 10)

	for _, p := range posts {
		if p.CreatedAt.Before(p.LastModified) {
			deletables = append(deletables, p)
		}
	}

	if len(deletables) > 0 {
		deleteRequest := buildDeleteRequest(deletables)

		for i := 1; ; i++ {
			pipeReader, pipeWriter := io.Pipe()

			go func() {
				jsonEncoder := json.NewEncoder(pipeWriter)
				err := jsonEncoder.Encode(&deleteRequest)
				pipeWriter.CloseWithError(err)
			}()

			r, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/post_%s/documents/query", s.host, board), pipeReader)
			resp, err := s.client.Do(r)

			if err != nil {
				if i < 3 {
					log.Printf("Error performing deletion request: %s", err)
					time.Sleep(30 * time.Second)
					continue
				} else {
					return fmt.Errorf("Error performing deletion request: %s", err)
				}
			}

			resp.Body.Close()

			if resp.StatusCode != 200 {
				return fmt.Errorf("Error deleting old posts: request received status %s", resp.Status)
			}

			break
		}
	}

	lnxPosts := DbPostsToLnxPosts(posts)

	for i := 1; true; i++ {
		pipeReader, pipeWriter := io.Pipe()

		go func() {
			jsonEncoder := json.NewEncoder(pipeWriter)
			err := jsonEncoder.Encode(&lnxPosts)
			pipeWriter.CloseWithError(err)
		}()

		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/documents", s.host, board), "application/json", pipeReader)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing insertion request: %s", err)
				time.Sleep(30 * time.Second)
				continue
			} else {
				return fmt.Errorf("Error performing insertion request: %s", err)
			}
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("Error inserting posts: request received status %s", resp.Status)
		}

		break
	}

	return nil
}

func (s *Service) Rollback(board string) error {
	for i := 1; ; i++ {
		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/rollback", s.host, board), "", nil)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing rollback: %s", err)
				continue
			} else {
				return fmt.Errorf("Error performing rollback: %s", err)
			}
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Rollback request received status %s", resp.Status)
		}

		return nil
	}
}

func (s *Service) Commit(board string) error {
	for i := 1; ; i++ {
		resp, err := s.client.Post(fmt.Sprintf("%s/post_%s/commit", s.host, board), "", nil)

		if err != nil {
			if i < 3 {
				log.Printf("Error performing commit: %s", err)
				continue
			} else {
				return fmt.Errorf("Error performing commit: %s", err)
			}
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Request received status %s", resp.Status)
		}

		return nil
	}
}
