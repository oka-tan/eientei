package lnx

import (
	"fmt"
	"moon/db"
	"strings"
)

func buildDeleteRequest(posts []db.Post) deleteRequest {
	var b strings.Builder

	fmt.Fprintf(&b, "no:%d", posts[0].No)
	for _, p := range posts[1:] {
		fmt.Fprintf(&b, " OR no:%d", p.No)
	}

	return deleteRequest{
		Query: query{normalQuery{Ctx: b.String()}},
		Limit: uint64(len(posts)),
	}
}
