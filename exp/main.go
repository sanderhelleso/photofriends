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

	type User struct {
		ID int
		Name string
		Email string
	}

	var users []User

	rows, err := db.Query(`
		SELECT id, name, email
		FROM users`)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			panic(err)
		}

		users = append(users, user)
	}
	fmt.Println(users)
}

/*CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT NOT NULL
);*/