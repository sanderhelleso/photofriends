package controllers

import (
	"../../photofriends/views"
	"net/http"
)

// NewUsers is uused to create a new Users controller
// this function will panic if the templates are not
// passed correctly, and should only be used during
// initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("layout", "views/users/new.gohtml"),
	}
}

// users struct containing the releated view
type Users struct {
	NewView *views.View
}

// method that automaticly pass in a Users pointer on call
func (u *Users) New(res http.ResponseWriter, req *http.Request) {
	if err := u.NewView.Render(res, nil); err != nil {
		panic(err)
	}
}