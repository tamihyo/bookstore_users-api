package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tamihyo/bookstore_users-api/app"
	"github.com/tamihyo/bookstore_users-api/datasources/mysql/users_db"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error getting env, not comming through %v", err)
	}
	users_db.InitDB()
	defer app.StartApplication()
}
