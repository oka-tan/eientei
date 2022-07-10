package db

import (
	"context"
	"kaguya/api"
	"kaguya/config"

	"github.com/uptrace/bun"
)

var batchSize int = 80

//Service encapsulates operations on the database.
type Service struct {
	boardName string
	db        *bun.DB
}

//NewService creates and returns a db.Service
func NewService(
	db *bun.DB,
	boardConfig config.BoardConfig,
) Service {
	return Service{
		db:        db,
		boardName: boardConfig.Name,
	}
}

//Insert inserts a slice of posts into the database.
//Behavior for conflicts is doing nothing.
//For updating on conflicts use Upsert.
func (s *Service) Insert(
	newPosts []api.Post,
) error {
	if len(newPosts) == 0 {
		return nil
	}

	posts := make([]post, 0, min(batchSize, len(newPosts)))

	for len(newPosts) > batchSize {
		for _, p := range newPosts[:batchSize] {
			posts = append(posts, toModel(s.boardName, p))
		}

		_, err := s.db.
			NewInsert().
			Model(&posts).
			On("CONFLICT (board, no) DO NOTHING").
			Returning("NULL").
			Exec(context.Background())

		if err != nil {
			return err
		}

		newPosts = newPosts[batchSize:]
		posts = posts[0:0]
	}

	for _, p := range newPosts {
		posts = append(posts, toModel(s.boardName, p))
	}

	_, err := s.db.
		NewInsert().
		Model(&posts).
		On("CONFLICT (board, no) DO NOTHING").
		Returning("NULL").
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

//Upsert inserts a slice of posts into the database and performs updates on conflict.
//For doing nothing on conflicts use Insert.
func (s *Service) Upsert(
	newPosts []api.Post,
) error {
	if len(newPosts) == 0 {
		return nil
	}

	posts := make([]post, 0, min(batchSize, len(newPosts)))

	for len(newPosts) > batchSize {
		for _, p := range newPosts[:batchSize] {
			posts = append(posts, toModel(s.boardName, p))
		}

		_, err := s.db.
			NewInsert().
			Model(&posts).
			On("CONFLICT (board, no) DO UPDATE").
			Set("com = EXCLUDED.com, file_deleted = EXCLUDED.file_deleted, sticky = EXCLUDED.sticky").
			Returning("NULL").
			Exec(context.Background())

		if err != nil {
			return err
		}

		newPosts = newPosts[batchSize:]
		posts = posts[0:0]
	}

	for _, p := range newPosts {
		posts = append(posts, toModel(s.boardName, p))
	}

	_, err := s.db.
		NewInsert().
		Model(&posts).
		On("CONFLICT (board, no) DO UPDATE").
		Set("com = EXCLUDED.com, file_deleted = EXCLUDED.file_deleted, sticky = EXCLUDED.sticky").
		Returning("NULL").
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

//Update updates a list of posts.
//For updating or inserting, use Upsert.
func (s *Service) Update(updatedPosts []api.Post) error {
	if len(updatedPosts) == 0 {
		return nil
	}

	posts := make([]post, 0, min(len(updatedPosts), batchSize))

	for len(updatedPosts) > batchSize {
		for _, p := range updatedPosts[:batchSize] {
			posts = append(posts, toModel(s.boardName, p))
		}

		_, err := s.db.NewUpdate().
			With("_data", s.db.NewValues(posts)).
			Model((*post)(nil)).
			TableExpr("_data").
			Set("com = _data.com").
			Set("file_deleted = _data.file_deleted").
			Set("sticky = _data.sticky").
			Where("post.no = _data.no").
			Returning("NULL").
			Exec(context.Background())

		if err != nil {
			return err
		}

		updatedPosts = updatedPosts[batchSize:]
		posts = posts[0:0]
	}

	for _, p := range updatedPosts {
		posts = append(posts, toModel(s.boardName, p))
	}

	_, err := s.db.NewUpdate().
		With("_data", s.db.NewValues(posts)).
		Model((*post)(nil)).
		TableExpr("_data").
		Set("com = _data.com").
		Set("file_deleted = _data.file_deleted").
		Set("sticky = _data.sticky").
		Where("post.no = _data.no").
		Returning("NULL").
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

//Delete marks a slice of posts as deleted.
func (s *Service) Delete(posts []int64, board string) error {
	if len(posts) == 0 {
		return nil
	}

	_, err := s.db.NewUpdate().
		Model((*post)(nil)).
		Set("deleted  = TRUE").
		Where("board = ?", board).
		Where("post.no IN (?)", bun.In(posts)).
		Returning("NULL").
		Exec(context.Background())

	return err
}
