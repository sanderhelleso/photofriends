package main

import (
	"fmt"
	"net/http"
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func handlerFunc(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")

	if req.URL.Path == "/" {
		fmt.Fprint(res, "<h1>Welcome to my awesome site!</h1>")
	} else if req.URL.Path == "/contact" {
		fmt.Fprint(res, "To get in touch, please send an email to <a href=\"mailto:support@photofriends.com\">support@photofriends.com</a>.")
	} else {
		fmt.Fprint(res, "<h1>404! Page not found</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc) // route and matched content
	http.ListenAndServe(":3000", nil) // port to serve (nil = NULLPOINTER)
}
