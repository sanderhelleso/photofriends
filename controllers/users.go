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
		LoginView: views.NewView("layout", "users/login"),
		us: us,
	}
}

type Users struct {
	NewView 	*views.View
	LoginView   *views.View
	us  		*models.UserService
}

type SignupForm struct {
	Name	 string	`schema:"name"`
	Email 	 string `schema:"email"`
	Password string `schema:"password"`
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

type LoginForm struct {
	Email 	 string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify the provided email address and password
// and then log the user in if the provided info is correct
//
// POST /login
func (u *Users) Login(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	var form LoginForm
	if err := dec.Decode(&form, req.PostForm); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(res, "Invalid email address.")
	case models.ErrInvalidPassword:
		fmt.Fprintln(res, "Invalid password provided.")
	case nil:
		fmt.Fprintln(res, user)
	default:
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}