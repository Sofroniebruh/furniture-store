package main

import (
	"auth-service/db"
	"auth-service/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	err := db.Init()

	if err != nil {
		log.Fatal("Error initializing database: ", err)
		return
	}

	r := chi.NewRouter()

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	r.Post("/registration", handlers.Signup)
	r.Post("/login", handlers.Login)
	r.Post("/logout", handlers.Logout)
	r.Post("/refresh", handlers.Refresh)

	log.Println("Listening on port 8081")
	err = http.ListenAndServe(":8081", handler)

	if err != nil {
		log.Fatal("Failed listening on port 8081: ", err)
		return
	}
}
