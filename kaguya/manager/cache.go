package manager

type cachedThread struct {
	lastModified uint64
	posts        map[int64]cachedPost
}

type cachedPost struct {
	comHash     uint32
	fileDeleted bool
	sticky      bool
}
