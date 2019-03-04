package main

import (
	"fmt"
	"net/http"
)

/*  Similar as node.js express (req, res)
Send back data to matched path using the writer(w)
*/
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc) // route and matched content
	http.ListenAndServe(":3000", nil) // port to serve
}
