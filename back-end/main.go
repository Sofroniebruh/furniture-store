package main

import (
	"furniture-store-backend/db"
	"furniture-store-backend/handlers"
	"log"
	"net/http"
)

func main() {
	err := db.Init()

	if err != nil {
		log.Fatal("Error initializing database: ", err)
		return
	}

	http.HandleFunc("/registration", handlers.Signup)
	http.HandleFunc("/login", handlers.Login)

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
