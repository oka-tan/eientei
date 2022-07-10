package main

import (
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	Board         string    `bun:"board" json:"board"`
	No            int64     `bun:"no,pk" json:"no"`
	Resto         int64     `bun:"resto,notnull" json:"resto"`
	Time          time.Time `bun:"time,notnull" json:"time"`
	Name          *string   `bun:"name" json:"name"`
	Trip          *string   `bun:"trip" json:"trip"`
	Capcode       *string   `bun:"capcode" json:"capcode"`
	Country       *string   `bun:"country" json:"country"`
	Since4Pass    *int16    `bun:"since4pass" json:"since4pass"`
	Sub           *string   `bun:"sub" json:"sub"`
	Com           *string   `bun:"com" json:"com"`
	Tim           *int64    `bun:"tim" json:"tim"`
	MD5           *string   `bun:"md5" json:"md5"`
	Filename      *string   `bun:"filename" json:"filename"`
	Ext           *string   `bun:"ext" json:"ext"`
	Fsize         *int64    `bun:"fsize" json:"fsize"`
	W             *int16    `bun:"w" json:"w"`
	H             *int16    `bun:"h" json:"h"`
	TnW           *int16    `bun:"tn_w" json:"tn_w"`
	TnH           *int16    `bun:"tn_h" json:"tn_h"`
	Deleted       bool      `bun:"deleted" json:"deleted"`
	FileDeleted   bool      `bun:"file_deleted" json:"file_deleted"`
	Spoiler       bool      `bun:"spoiler" json:"spoiler"`
	CustomSpoiler *int8     `bun:"custom_spoiler" json:"custom_spoiler"`
	Op            bool      `bun:"op" json:"op"`
	Sticky        bool      `bun:"sticky" json:"sticky"`
}
