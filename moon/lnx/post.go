package lnx

import (
	"moon/db"
	"time"
)

type Post struct {
	No          int64     `json:"no"`
	Resto       int64     `json:"resto"`
	Time        time.Time `json:"time"`
	Name        string    `json:"name"`
	Trip        string    `json:"trip"`
	Capcode     string    `json:"capcode"`
	Country     string    `json:"country"`
	Since4Pass  int16     `json:"since4pass"`
	Sub         string    `json:"sub"`
	Com         string    `json:"com"`
	Tim         int64     `json:"tim"`
	MD5         string    `json:"md5"`
	Filename    string    `json:"filename"`
	Ext         string    `json:"ext"`
	Deleted     int8      `json:"deleted"`
	FileDeleted int8      `json:"file_deleted"`
	Spoiler     int8      `json:"spoiler"`
	Op          int8      `json:"op"`
	Sticky      int8      `json:"sticky"`
}

func DbPostToLnxPost(p db.Post) Post {
	return Post{
		No:          p.No,
		Resto:       p.Resto,
		Time:        p.Time,
		Name:        derefString(p.Name),
		Trip:        derefString(p.Trip),
		Capcode:     derefString(p.Capcode),
		Country:     derefString(p.Country),
		Since4Pass:  derefInt16(p.Since4Pass),
		Sub:         derefString(p.Sub),
		Com:         derefString(p.Com),
		Tim:         derefInt64(p.Tim),
		MD5:         derefString(p.MD5),
		Filename:    derefString(p.Filename),
		Ext:         derefString(p.Ext),
		Deleted:     boolToInt8(p.Deleted),
		FileDeleted: boolToInt8(p.FileDeleted),
		Spoiler:     boolToInt8(p.Spoiler),
		Op:          boolToInt8(p.Op),
		Sticky:      boolToInt8(p.Sticky),
	}
}

func DbPostsToLnxPosts(posts []db.Post) []Post {
	lnxPosts := make([]Post, 0, len(posts))

	for _, p := range posts {
		lnxPosts = append(lnxPosts, DbPostToLnxPost(p))
	}

	return lnxPosts
}

func boolToInt8(b bool) int8 {
	if b {
		return 1
	} else {
		return 0
	}
}

func derefInt16(i *int16) int16 {
	if i != nil {
		return *i
	} else {
		return 0
	}
}

func derefInt64(i *int64) int64 {
	if i != nil {
		return *i
	} else {
		return 0
	}
}

func derefString(s *string) string {
	if s != nil {
		return *s
	} else {
		return ""
	}
}
