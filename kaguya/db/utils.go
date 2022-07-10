package db

import (
	"kaguya/api"
	"kaguya/utils"
	"time"
)

func filterName(name *string) *string {
	if name == nil || *name == "Anonymous" {
		return nil
	}

	return name
}

func toModel(boardName string, p api.Post) post {
	var resto int64
	if p.Resto != 0 {
		resto = p.Resto
	} else {
		resto = p.No
	}

	return post{
		Board:         boardName,
		No:            p.No,
		Resto:         resto,
		Time:          time.UnixMilli(1000 * int64(p.Time)),
		Name:          filterName(p.Name),
		Trip:          p.Trip,
		Capcode:       p.Capcode,
		Country:       p.Country,
		Since4Pass:    p.Since4Pass,
		Sub:           p.Sub,
		Com:           p.Com,
		Tim:           p.Tim,
		MD5:           p.MD5,
		Filename:      p.Filename,
		Ext:           p.Ext,
		Fsize:         p.Fsize,
		W:             p.W,
		H:             p.H,
		TnW:           p.TnW,
		TnH:           p.TnH,
		Deleted:       false,
		FileDeleted:   utils.ToBool(p.FileDeleted),
		Spoiler:       utils.ToBool(p.Spoiler),
		CustomSpoiler: p.CustomSpoiler,
		Op:            (p.Resto == 0),
		Sticky:        utils.ToBool(p.Sticky),
	}
}

func min(i1 int, i2 int) int {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}
