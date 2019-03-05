package main

import (
	"fmt"
	"../../photofriends/models"
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

	us, err := models.NewUserService(psqlInfo); if err != nil {
		panic(err);
	}

	defer us.Close()
	us.DestructiveReset()

	//user, err := us.ByID(1)
	//fmt.Println(user)
}