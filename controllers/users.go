package controllers

import (
	"fmt"
	"net/http"

	"../../photofriends/models"
	"../../photofriends/rand"
	"../../photofriends/views"
	"github.com/gorilla/schema"
)

// NewUsers is uused to create a new Users controller
// this function will panic if the templates are not
// passed correctly, and should only be used during
// initial setup
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("layout", "users/new"),
		LoginView: views.NewView("layout", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
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

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err := u.signIn(res, &user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

type LoginForm struct {
	Email    string `schema:"email"`
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
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(res, "Invalid email address.")
		case models.ErrPasswordIncorrect:
			fmt.Fprintln(res, "Invalid password provided.")
		default:
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	}

	err = u.signIn(res, user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(res, req, "/cookietest", http.StatusFound)
}

// signIn is used to sign in the given user in via cookies
func (u *Users) signIn(res http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token
		err = u.us.Update(user)
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(res, &cookie)
	return nil
}

// CookieTest is used to display cookies on the current user
func (u *Users) CookieTest(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("remember_token")
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusFound)
		return
	}
	fmt.Fprintln(res, user)
}
