package main

import (
	"fmt"
	"net/http"

	"../photofriends/controllers"
	"../photofriends/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "photofriends_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	must(err)

	defer services.Close()
	services.DestructiveReset()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

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
	router.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// gallery routes
	router.Handle("/galleries/new", galleriesC.New).Methods("GET")
	router.HandleFunc("/galleries", galleriesC.Create).Methods("POST")

	http.ListenAndServe(":3000", router) // port to serve (nil = NULLPOINTER)
}

// panic if ANY error is present
func must(err error) {
	if err != nil {
		panic(err)
	}
}
