package controllers

import (
	"../../photofriends/views"
	"../../photofriends/models"
	"github.com/gorilla/schema"
	"net/http"
	"fmt"
)

// NewUsers is uused to create a new Users controller
// this function will panic if the templates are not
// passed correctly, and should only be used during
// initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("layout", "users/new"),
		us: us,
	}
}

type Users struct {
	NewView *views.View
	us  *models.UserService
}

type SignupForm struct {
	Name	 string	`schema:"name"`
	Email 	 string `schema:"email"`
	Password string `schema:"password"`
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

	dec := schema.NewDecoder()
	var form SignupForm
	if err := dec.Decode(&form, req.PostForm); err != nil {
		panic(err)
	}

	user := models.User {
		Name : form.Name,
		Email: form.Email,
		Password: form.Password,
	}
	
	if err := u.us.Create(&user); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(res, user)
}