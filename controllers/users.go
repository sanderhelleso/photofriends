package controllers

import (
	"../../photofriends/views"
	"net/http"
	"fmt"
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

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(res http.ResponseWriter, req *http.Request) {
	if err := u.NewView.Render(res, nil); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form when a user
// submits the form. This is used to create a new user account
//
// POST /signup
func (u *Users) Create(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	fmt.Fprintln(res, req.PostForm["email"])
	fmt.Fprintln(res, req.PostFormValue("email"))
	fmt.Fprintln(res, req.PostForm["password"])
	fmt.Fprintln(res, req.PostFormValue("password"))
	fmt.Fprintln(res, "This is a temp response")
}