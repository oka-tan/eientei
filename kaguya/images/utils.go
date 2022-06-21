package images

import "kaguya/api"

func toQueuedImages(boardName string, posts []api.Post) []QueuedImage {
	batch := make([]QueuedImage, 0, 10)

	for _, p := range posts {
		if p.Filename != nil && p.Tim != nil && p.Fsize != nil && p.Ext != nil {
			batch = append(batch, QueuedImage{
				Tim:   *p.Tim,
				Fsize: *p.Fsize,
				Ext:   *p.Ext,
				Board: boardName,
			})
		}
	}

	return batch
}
