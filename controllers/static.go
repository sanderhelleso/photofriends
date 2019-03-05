package controllers

import "../../photofriends/views"

func NewStatic() *Static {
	return &Static {
		Home: 	 views.NewView("layout", "views/static/home.gohtml"),
		Contact: views.NewView("layout", "views/static/contact.gohtml"),
	}
}

type Static struct {
	Home	 *views.View
	Contact  *views.View
}