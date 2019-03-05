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
	must(homeView.Render(res, nil))
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	must(contactView.Render(res, nil))
}

func faq(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, "Frequently asked questions: Comming Soon")
}

func notFound(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "text/html")
	res.WriteHeader(http.StatusNotFound)
	fmt.Fprint(res, "<p>404 page not found</p>")
}

func main() {

	// view setup
	homeView = views.NewView("layout", "views/home.gohtml")
	contactView = views.NewView("layout", "views/contact.gohtml")

	// router & path config
	router := mux.NewRouter()                           // router
	router.NotFoundHandler = http.HandlerFunc(notFound) // 404 not found
	router.HandleFunc("/", home)                        // home
	router.HandleFunc("/contact", contact)              // contact
	router.HandleFunc("/fag", faq)                      // faq
	http.ListenAndServe(":3000", router)                // port to serve (nil = NULLPOINTER)
}

// panic if ANY error is present
func must(err error) {
	if err != nil {
		panic(err);
	}
}
