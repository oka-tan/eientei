package manager

import (
	"hash/fnv"
	"kaguya/api"
	"kaguya/utils"
)

func comHash(com *string) uint32 {
	var comHash uint32

	if com != nil {
		hash := fnv.New32()
		hash.Sum([]byte(*com))
		comHash = hash.Sum32()
	} else {
		comHash = 0
	}

	return comHash
}

func toCachedPosts(thread map[int64]api.Post) map[int64]cachedPost {
	postsCache := make(map[int64]cachedPost)

	for _, p := range thread {
		postsCache[p.No] = toCachedPost(p)
	}

	return postsCache
}

func toCachedPost(p api.Post) cachedPost {
	return cachedPost{
		comHash:     comHash(p.Com),
		fileDeleted: utils.ToBool(p.FileDeleted),
		sticky:      utils.ToBool(p.Sticky),
	}
}

func modified(p1 api.Post, p2 cachedPost) bool {
	return utils.ToBool(p1.Sticky) != p2.sticky || utils.ToBool(p1.FileDeleted) != p2.fileDeleted || comHash(p1.Com) != p2.comHash
}
