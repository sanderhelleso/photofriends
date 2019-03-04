package main

import (
	"fmt"
	"net/http"
)

// Similar as node.js express (req, res)
// Send back data to matched path using the writer(res)
func handlerFunc(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	fmt.Fprint(res, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc) // route and matched content
	http.ListenAndServe(":3000", nil) // port to serve (nil = NULLPOINTER)
}
