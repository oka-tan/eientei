package manager

import (
	"fmt"
	"kaguya/api"
	"kaguya/config"
	"kaguya/db"
	"kaguya/images"
	"kaguya/thumbnails"
	"log"
	"time"

	"github.com/uptrace/bun"
)

//BoardManager is a service responsible for organizing and timing the scraping.
type BoardManager struct {
	apiService        api.Service
	dbService         db.Service
	imagesService     images.Service
	thumbnailsService thumbnails.Service
	boardName         string
	threads           map[int64]cachedThread
	napTime           time.Duration
	longNapTime       time.Duration
	images            bool
	thumbnails        bool
}

//NewBoardManager creates and returns a BoardManager instance specific to some board.
//BoardManager instances need to be Init'ed and then should be left to Run.
func NewBoardManager(
	apiConfig config.APIConfig,
	boardConfig config.BoardConfig,
	pg *bun.DB,
	imagesService images.Service,
	thumbnailsService thumbnails.Service,
) (BoardManager, error) {
	apiService, err := api.NewService(
		apiConfig,
		boardConfig,
	)

	if err != nil {
		return BoardManager{}, fmt.Errorf(
			"Error constructing API interface:\nAPI Configuration %+v\nBoard Configuration:%+v",
			apiConfig,
			boardConfig,
		)
	}

	dbService := db.NewService(
		pg,
		boardConfig,
	)

	threads := make(map[int64]cachedThread)

	napTime, _ := time.ParseDuration(boardConfig.NapTime)
	longNapTime, _ := time.ParseDuration(boardConfig.LongNapTime)

	return BoardManager{
		apiService:        apiService,
		dbService:         dbService,
		boardName:         boardConfig.Name,
		imagesService:     imagesService,
		thumbnailsService: thumbnailsService,
		threads:           threads,
		napTime:           napTime,
		longNapTime:       longNapTime,
		images:            boardConfig.Images,
		thumbnails:        boardConfig.Thumbnails,
	}, nil
}

//Init puts a BoardManager in working conditions.
func (c *BoardManager) Init() error {
	catalog, err := c.apiService.GetCatalog()

	if err != nil {
		return err
	}

	archive, err := c.apiService.GetArchive()

	if err != nil {
		return err
	}

	for no := range archive {
		time.Sleep(c.napTime)

		posts, err := c.apiService.GetThreadArray(no)

		if err != nil {
			log.Printf("Error looking up thread %d: %s\n", no, err)
			continue
		}

		if err = c.dbService.Upsert(posts); err != nil {
			return err
		}

		if c.images {
			c.imagesService.Enqueue(c.boardName, posts)
		}

		if c.thumbnails {
			c.thumbnailsService.Enqueue(c.boardName, posts)
		}
	}

	log.Printf("Finished up with the archive for board %s, moving on to the catalog\n", c.boardName)

	for _, thread := range catalog {
		time.Sleep(c.napTime)

		posts, err := c.apiService.GetThreadArray(thread.No)

		if err != nil {
			log.Printf("Error looking up thread %d: %s", thread.No, err)
			continue
		}

		if err = c.dbService.Upsert(posts); err != nil {
			return err
		}

		if c.images {
			c.imagesService.Enqueue(c.boardName, posts)
		}

		if c.thumbnails {
			c.thumbnailsService.Enqueue(c.boardName, posts)
		}

		postsMap := make(map[int64]api.Post)

		for _, p := range posts {
			postsMap[p.No] = p
		}

		c.threads[thread.No] = cachedThread{
			lastModified: thread.LastModified,
			posts:        toCachedPosts(postsMap),
		}
	}

	return nil
}

//Run is the main method for a BoardManager.
//Infinite loop that scrapes and repeats.
func (c *BoardManager) Run() error {
	for {
		time.Sleep(c.longNapTime)
		log.Printf("Starting loop for board %s", c.boardName)

		catalog, err := c.apiService.GetCatalog()

		if err != nil {
			log.Printf("Error looking up catalog: %s", err)
			continue
		}

		time.Sleep(c.napTime)
		archive, err := c.apiService.GetArchive()

		if err != nil {
			log.Printf("Error looking up archive: %s", err)
			archive = make(map[int64]bool)
		}

		deletedPosts := make([]int64, 0, 10)
		editedPosts := make([]api.Post, 0, 10)
		newPosts := make([]api.Post, 0, 10)

		for no := range c.threads {
			_, inCatalog := catalog[no]
			_, inArchive := archive[no]

			if !inCatalog && !inArchive {
				log.Printf("Marking thread %d as deleted", no)

				deletedPosts = append(deletedPosts, no)
				delete(c.threads, no)
			}
		}

		for threadNo := range archive {
			oldThread, threadExisted := c.threads[threadNo]

			if !threadExisted {
				continue
			}

			time.Sleep(c.napTime)
			updatedThread, err := c.apiService.GetThread(threadNo)

			if err != nil {
				log.Printf("Error loading archived thread %d: %s", threadNo, err)
				delete(c.threads, threadNo)
				continue
			}

			for no, p := range oldThread.posts {
				updatedPost, notDeleted := updatedThread[no]

				if !notDeleted {
					deletedPosts = append(deletedPosts, no)
					continue
				}

				if modified(updatedPost, p) {
					editedPosts = append(editedPosts, updatedPost)
					delete(updatedThread, no)
				}
			}

			for _, p := range updatedThread {
				newPosts = append(newPosts, p)
			}

			delete(c.threads, threadNo)
		}

		for threadNo, t := range catalog {
			oldThread, threadExisted := c.threads[threadNo]

			if threadExisted && oldThread.lastModified == t.LastModified {
				continue
			}

			time.Sleep(c.napTime)

			updatedThread, err := c.apiService.GetThread(threadNo)

			if err != nil {
				log.Printf("Error loading thread %d: %s", threadNo, err)
				continue
			}

			if !threadExisted {
				for _, p := range updatedThread {
					newPosts = append(newPosts, p)
				}

				c.threads[threadNo] = cachedThread{
					lastModified: t.LastModified,
					posts:        toCachedPosts(updatedThread),
				}

				continue
			}

			oldPosts := oldThread.posts

			for no, p := range oldPosts {
				updatedPost, notDeleted := updatedThread[no]

				if !notDeleted {
					deletedPosts = append(deletedPosts, no)
					delete(oldPosts, no)

					continue
				}

				if modified(updatedPost, p) {
					editedPosts = append(editedPosts, updatedPost)
					delete(updatedThread, no)
					oldPosts[no] = toCachedPost(updatedPost)
				}
			}

			for no, p := range updatedThread {
				oldPosts[no] = toCachedPost(p)
				newPosts = append(newPosts, p)
			}

			c.threads[threadNo] = cachedThread{
				lastModified: t.LastModified,
				posts:        oldPosts,
			}
		}

		if err = c.dbService.Delete(deletedPosts); err != nil {
			return err
		}

		if c.dbService.Update(editedPosts); err != nil {
			return err
		}

		if c.images {
			c.imagesService.Enqueue(c.boardName, newPosts)
		}

		if c.thumbnails {
			c.thumbnailsService.Enqueue(c.boardName, newPosts)
		}

		if c.dbService.Insert(newPosts); err != nil {
			return err
		}
	}
}
