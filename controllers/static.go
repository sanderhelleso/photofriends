package controllers

import "../../photofriends/views"

func NewStatic() *Static {
	return &Static {
		Home: 	 views.NewView("layout", "static/home"),
		Contact: views.NewView("layout", "static/contact"),
	}
}

type Static struct {
	Home	 *views.View
	Contact  *views.View
}