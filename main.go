package main

import (
	"net/http"

	"../photofriends/controllers"

	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	usersC  := controllers.NewUsers()

	// router & path config
	// note the "Methods", it specify that 
	// only the sat requests types are allowed
	router := mux.NewRouter() // router
	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET") 
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", router) // port to serve (nil = NULLPOINTER)
}

// panic if ANY error is present
func must(err error) {
	if err != nil {
		panic(err);
	}
}
