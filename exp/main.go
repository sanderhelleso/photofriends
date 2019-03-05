package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil { panic(err) }
	defer db.Close()

	var id int
	err = db.QueryRow(`
		INSERT INTO users(name, email)
		VALUES ($1, $2)
		RETURNING id`, 
		"John Doe", "johndoe@gmail.com").Scan(&id)

	if err != nil {
		panic(err)
	}

	fmt.Println("new id is...", id)
}

/*CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT NOT NULL
);*/