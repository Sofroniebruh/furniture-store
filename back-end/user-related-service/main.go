package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
	"user-related-service/db"
	"user-related-service/middleware"
	"user-related-service/services"
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

	r.Get("/user", middleware.Protected(http.HandlerFunc(services.GetUserInfo)))

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", handler)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
