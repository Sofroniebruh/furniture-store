package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"log"
	"net/http"
	"products-service/db"
	"products-service/middleware"
	"products-service/services"
)

func main() {
	err := db.Init()

	if err != nil {
		log.Fatal("Error initializing database: ", err)
		return
	}

	r := chi.NewRouter()

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://fumi.artorien.me"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.ProtectedWithRoles("admin"))
		r.Put("/products", services.UpdateProduct)
		r.Delete("/products", services.DeleteProduct)
		r.Post("/products", services.AddProduct)

		r.Post("/colors", services.CreateColor)
		r.Delete("/colors", services.DeleteColor)
		r.Put("/colors", services.UpdateColor)

		r.Post("/stock", services.UpdateProductStock)
		r.Get("/stock-history", services.GetStockHistory)
		r.Get("/products/{id}/stock-history", services.GetProductStockHistory)
	})

	r.Get("/products", services.GetProducts)
	r.Get("/products/{id}", services.GetProductById)

	r.Get("/colors", services.GetAllColors)

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(":8080", handler)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
