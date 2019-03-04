package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, "<h1>Welcome to my awesome site!</h1>")
}

func contact(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, "To get in touch, please send an email to <a href=\"mailto:support@photofriends.com\">support@photofriends.com</a>.")
}

func main() {
	router := mux.NewRouter()              // router
	router.HandleFunc("/", home)           // home
	router.HandleFunc("/contact", contact) // contact
	http.ListenAndServe(":3000", router)   // port to serve (nil = NULLPOINTER)
}
