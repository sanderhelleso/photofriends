package controllers

import (
	"../../photofriends/models"
	"../../photofriends/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("layout", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
