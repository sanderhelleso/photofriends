package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
  	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host 	 = "localhost"
	port 	 = 5432
	user	 = "postgres"
	password = "postgres"

	dbname 	 = "photofriends_dev"
)

type User struct {
	gorm.Model
	Name string
	Email string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	
	if err != nil { panic(err) }

	defer db.Close()
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	//db.DropTableIfExists($User{})
	//db.AutoMigrate(&User{}) // automigrates the user struct as table
}

/*CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT NOT NULL
);*/