package controllers

import (
	"fmt"
	"net/http"

	"../../photofriends/models"
	"../../photofriends/views"
	"../context"
	"github.com/gorilla/schema"
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

type GalleryForm struct {
	Title string `schema:"title"`
}

// POST /galleries
func (g *Galleries) Create(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	var form GalleryForm
	if err := dec.Decode(&form, req.PostForm); err != nil {
		return
	}

	user := context.User(req.Context())
	if user == nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(res, gallery)
}
