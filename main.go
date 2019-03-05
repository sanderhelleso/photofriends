package main

import (
	"fmt"
	"net/http"

	"../photofriends/views"

	"github.com/gorilla/mux"
)

// views initializing
var (
	homeView    *views.View
	contactView *views.View
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	err := homeView.Template.Execute(res, nil)
	if err != nil { panic(err) }
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	err := contactView.Template.Execute(res, nil)
	if err != nil { panic(err) }
}

func faq(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, "Frequently asked questions: Comming Soon")
}

func notFound(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res, "<p>404 page not found")
}

func main() {

	// view setup
	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	// router & path config
	router := mux.NewRouter()                           // router
	router.NotFoundHandler = http.HandlerFunc(notFound) // 404 not found
	router.HandleFunc("/", home)                        // home
	router.HandleFunc("/contact", contact)              // contact
	router.HandleFunc("/fag", faq)                      // faq
	http.ListenAndServe(":3000", router)                // port to serve (nil = NULLPOINTER)
}
