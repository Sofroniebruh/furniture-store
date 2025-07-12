package main

import (
	"furniture-store-backend/db"
	"furniture-store-backend/handlers"
	"furniture-store-backend/middleware"
	"furniture-store-backend/services"
	"github.com/go-chi/chi/v5"
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

	r.Post("/registration", handlers.Signup)
	r.Post("/login", handlers.Login)
	r.Post("/logout", handlers.Logout)
	r.Post("/refresh", handlers.Refresh)
	r.Get("/products", services.GetProducts)

	r.Post("/products", middleware.Protected(http.HandlerFunc(services.AddProduct)))

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
