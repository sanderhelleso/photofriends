package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// template initializing
var (
	homeTemplate    *template.Template
	contactTemplate *template.Template
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(res, nil); err != nil {
		panic(err)
	}
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(res, nil); err != nil {
		panic(err)
	}
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

	// template setup
	var err error
	homeTemplate, err = template.ParseFiles(
		"views/home.gohtml",
		"views/layouts/footer.gohtml",
	)
	if err != nil {
		panic(err)
	}

	contactTemplate, err = template.ParseFiles(
		"views/contact.gohtml",
		"views/layouts/footer.gohtml",
	)
	if err != nil {
		panic(err)
	}

	// router & path config
	router := mux.NewRouter()                           // router
	router.NotFoundHandler = http.HandlerFunc(notFound) // 404 not found
	router.HandleFunc("/", home)                        // home
	router.HandleFunc("/contact", contact)              // contact
	router.HandleFunc("/fag", faq)                      // faq
	http.ListenAndServe(":3000", router)                // port to serve (nil = NULLPOINTER)
}
