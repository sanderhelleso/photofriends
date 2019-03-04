package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func handlerFunc(res http.ResponseWriter, req *http.Request) {

	// set content type to html
	res.Header().Set("Content-Type", "text/html")

	if req.URL.Path == "/" { // home
		fmt.Fprint(res, "<h1>Welcome to my awesome site!</h1>")
	} else if req.URL.Path == "/contact" { // contact
		fmt.Fprint(res, "To get in touch, please send an email to <a href=\"mailto:support@photofriends.com\">support@photofriends.com</a>.")
	} else { // 404
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "<h1>404! Page not found</h1>")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", r) // port to serve (nil = NULLPOINTER)
}
