package manager

type cachedThread struct {
	lastModified uint64
	posts        map[int64]cachedPost
}

type cachedPost struct {
	no          int64
	comHash     uint64
	fileDeleted bool
	sticky      bool
}
