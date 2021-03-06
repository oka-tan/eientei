package db

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:post,alias:post"`

	Board         string    `bun:"board"`
	No            int64     `bun:"no,pk"`
	Resto         int64     `bun:"resto,notnull"`
	Time          time.Time `bun:"time,notnull"`
	Name          *string   `bun:"name"`
	Trip          *string   `bun:"trip"`
	Capcode       *string   `bun:"capcode"`
	Country       *string   `bun:"country"`
	Since4Pass    *int16    `bun:"since4pass"`
	Sub           *string   `bun:"sub"`
	Com           *string   `bun:"com"`
	Tim           *int64    `bun:"tim"`
	MD5           *string   `bun:"md5"`
	Filename      *string   `bun:"filename"`
	Ext           *string   `bun:"ext"`
	Fsize         *int64    `bun:"fsize"`
	W             *int16    `bun:"w"`
	H             *int16    `bun:"h"`
	TnW           *int16    `bun:"tn_w"`
	TnH           *int16    `bun:"tn_h"`
	Deleted       bool      `bun:"deleted"`
	FileDeleted   bool      `bun:"file_deleted"`
	Spoiler       bool      `bun:"spoiler"`
	CustomSpoiler *int8     `bun:"custom_spoiler"`
	Op            bool      `bun:"op"`
	Sticky        bool      `bun:"sticky"`
}

func (p *Post) FormatName() string {
	if p.Name != nil {
		return *p.Name
	} else {
		return "Anonymous"
	}
}

func (p *Post) FormatTime() string {
	return p.Time.Format("Mon 2 Jan 2006 15:04:05")
}

func (p *Post) SubIsNil() bool {
	return p.Sub == nil
}

func (p *Post) DerefSub() string {
	return *(p.Sub)
}

func (p *Post) ComIsNil() bool {
	return p.Com == nil
}

func (p *Post) DerefCom() string {
	return *(p.Com)
}

func (p *Post) DerefTnH() int16 {
	return *(p.TnH)
}

func (p *Post) DerefTnW() int16 {
	return *(p.TnW)
}

func (p *Post) DerefH() int16 {
	return *(p.H)
}

func (p *Post) DerefW() int16 {
	return *(p.W)
}

func (p *Post) DerefFilename() string {
	return *(p.Filename)
}

func (p *Post) DerefFilenameShort() string {
	s := *(p.Filename)
	if len(s) < 20 {
		return s
	} else {
		runes := []rune(s)
		runesLen := len(runes)

		if runesLen > 17 {
			return fmt.Sprintf("%s...", string(runes[:17]))
		} else {
			return s
		}
	}
}

func (p *Post) TimIsNil() bool {
	return p.Tim == nil
}

func (p *Post) DerefTim() int64 {
	return *(p.Tim)
}

func (p *Post) DerefExt() string {
	return *(p.Ext)
}
