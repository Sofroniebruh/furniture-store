package main

import (
	"furniture-store-backend/db"
	"furniture-store-backend/handlers"
	"furniture-store-backend/middleware"
	"furniture-store-backend/services"
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

	r.Get("/products", services.GetProducts)
	r.Put("/products", middleware.Protected(http.HandlerFunc(services.UpdateProduct)))
	r.Delete("/products", middleware.Protected(http.HandlerFunc(services.DeleteProduct)))
	r.Post("/products", middleware.Protected(http.HandlerFunc(services.AddProduct)))

	r.Post("/colors", middleware.Protected(http.HandlerFunc(services.CreateColor)))
	r.Get("/colors", services.GetAllColors)
	r.Delete("/colors", middleware.Protected(http.HandlerFunc(services.DeleteColor)))
	r.Put("/colors", middleware.Protected(http.HandlerFunc(services.UpdateColor)))

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", handler)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
