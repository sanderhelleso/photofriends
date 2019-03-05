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

	user := models.User {
		Name: "Michael Scott",
		Email: "michael@dundermifflin.com",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	if err := us.Delete(user.ID); err != nil {
		panic(err)
	}


	userByID, err := us.ByID(user.ID)
	if err != nil { panic(err) }

	fmt.Println(userByID)
}