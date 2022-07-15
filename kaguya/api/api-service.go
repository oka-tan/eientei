package api

import (
	"encoding/json"
	"fmt"
	"kaguya/config"
	"net/http"
	"time"
)

//Service acts as an almost-singleton for consuming the 4chan API
type Service struct {
	boardName          string
	host               string
	client             http.Client
	lastCatalogRequest *time.Time
	lastArchiveRequest *time.Time
}

//NewService constructs an api.Service given the APIConfiguration and the configuration for some board
func NewService(
	apiConfig config.APIConfig,
	boardConfig config.BoardConfig,
) (Service, error) {
	requestTimeout, _ := time.ParseDuration(apiConfig.RequestTimeout)

	client := http.Client{
		Timeout: requestTimeout,
	}

	return Service{
		boardName: boardConfig.Name,
		host:      apiConfig.Host,
		client:    client,
	}, nil
}

//GetCatalog queries and returns the catalog as a map[int64]Catalog thread, where keys are thread numbers.
//If the catalog hasn't been modified since the last request returns nil on both arguments.
func (s *Service) GetCatalog() (map[int64]CatalogThread, error) {

	url := fmt.Sprintf("%s/%s/catalog.json", s.host, s.boardName)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error constructing catalog request:\nURL: %s\nError: %s", url, err)
	}

	if s.lastCatalogRequest != nil {
		req.Header["If-Modified-Since"] = []string{
			fmt.Sprintf("%s GMT", s.lastCatalogRequest.Format("Mon Mon, 02 Jan 2006 15:04:05")),
		}
	}

	utc, _ := time.LoadLocation("UTC")
	lastCatalogRequest := time.Now().In(utc)
	s.lastCatalogRequest = &lastCatalogRequest

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error performing catalog request:\nURL: %s\nError: %s", url, err)
	}

	defer resp.Body.Close()

	//Catalog hasn't been modified since the time specified in If-Modified-Since
	if resp.StatusCode == 304 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error processing catalog response:\nURL:%s\nError: status code is %d", url, resp.StatusCode)
	}

	var pages []CatalogPage

	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&pages)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response body of catalog request:\nURL: %s\nError: %s", url, err)
	}

	catalog := make(map[int64]CatalogThread)

	for _, page := range pages {
		for _, t := range page.Threads {
			catalog[t.No] = t
		}
	}

	return catalog, nil
}

//GetArchive queries and returns the board archive as a map[int64]bool (keys are thread numbers, bools are meaningless).
//If it hasn't been modified since the last request the archive is empty.
func (s *Service) GetArchive() (map[int64]bool, error) {

	url := fmt.Sprintf("%s/%s/archive.json", s.host, s.boardName)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error constructing archive request:\nURL: %s\nError: %s", url, err)
	}

	if s.lastArchiveRequest != nil {
		req.Header["If-Modified-Since"] = []string{
			fmt.Sprintf("%s GMT", s.lastArchiveRequest.Format("Mon Mon, 02 Jan 2006 15:04:05")),
		}
	}

	utc, _ := time.LoadLocation("UTC")
	lastArchiveRequest := time.Now().In(utc)
	s.lastArchiveRequest = &lastArchiveRequest

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error performing archive request:\nURL: %s\nError: %s", url, err)
	}

	defer resp.Body.Close()

	//Archive hasn't been modified since the time specified in If-Modified-Since
	if resp.StatusCode == 304 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error processing archive response:\nURL: %s\nError: status code is %d", url, resp.StatusCode)
	}

	var archiveJSON []int64

	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&archiveJSON)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling archive response body: %s", err)
	}

	archive := make(map[int64]bool)

	for _, no := range archiveJSON {
		archive[no] = true
	}

	return archive, nil
}

//GetThread looks up a thread and returns it as a map[int64]Post, where the map keys are the post numbers.
//
//The OP can be found by doing something like:
//for k, v := range t {
//  if v.ResTo == 0 {
//    return v
//  } else {
//    return t[v.ResTo]
//  }
//}
func (s *Service) GetThread(no int64) (map[int64]Post, error) {

	url := fmt.Sprintf("%s/%s/thread/%d.json", s.host, s.boardName, no)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error constructing thread request:\nURL: %s\nError: %s", url, err)
	}

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error performing thread request:\nURL: %s\nError: %s", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error processing thread response:\nURL: %s\nError: status code is %d", url, resp.StatusCode)
	}

	var thread Thread

	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&thread)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response body:\nError: %s", err)
	}

	posts := make(map[int64]Post)

	for _, p := range thread.Posts {
		posts[p.No] = p
	}

	return posts, nil
}

func (s *Service) GetThreadArray(no int64) ([]Post, error) {

	url := fmt.Sprintf("%s/%s/thread/%d.json", s.host, s.boardName, no)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error constructing thread request:\nURL: %s\nError: %s", url, err)
	}

	resp, err := s.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Error performing thread request:\nURL: %s\nError: %s", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error processing thread response:\nURL: %s\nError: status code is %d", url, resp.StatusCode)
	}

	var thread Thread

	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&thread)

	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling response body: %s", err)
	}

	return thread.Posts, nil
}
