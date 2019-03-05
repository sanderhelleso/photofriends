package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"

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
	gorm.Model // id, created at, deleted at, updated at
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
	db.LogMode(true)
	db.AutoMigrate(&User{}) // automigrates the user struct as table

	name, email := getInfo()
	user := User {
		Name: name,
		Email: email,
	}

	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", user)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	return;
}

/*CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT NOT NULL
);*/