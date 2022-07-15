package db

import (
	"time"

	"github.com/uptrace/bun"
)

type IndexTracker struct {
	bun.BaseModel `bun:"table:index_tracker"`

	Board        string    `bun:"board,pk"`
	LastModified time.Time `bun:"last_modified"`
	No           int64     `bun:"no"`
}
