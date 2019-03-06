package main

import (
	"net/http"
	"fmt"
	"../photofriends/models"
	"../photofriends/controllers"

	"github.com/gorilla/mux"
)

const (
	host 	 = "localhost"
	port 	 = 5432
	user	 = "postgres"
	password = "postgres"

	dbname 	 = "photofriends_dev"
)


func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	must(err)

	defer us.Close()
	//us.DestructiveReset()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC  := controllers.NewUsers(us)

	// router & path config
	// note the "Methods", it specify that 
	// only the sat requests types are allowed
	router := mux.NewRouter() // router
	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.Handle("/signup", usersC.NewView).Methods("GET") 
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	router.Handle("/login", usersC.LoginView).Methods("GET") 
	router.HandleFunc("/login", usersC.Login).Methods("POST")
	http.ListenAndServe(":3000", router) // port to serve (nil = NULLPOINTER)
}

// panic if ANY error is present
func must(err error) {
	if err != nil {
		panic(err);
	}
}
