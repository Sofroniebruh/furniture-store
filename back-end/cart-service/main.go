package main

import (
	"cart-service/config"
	"cart-service/db"
	"cart-service/handlers"
	authmiddleware "cart-service/middleware"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	err := db.Init()

	if err != nil {
		log.Fatal("Error initializing database: ", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://fumi.artorien.me"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	r.Post("/webhook", handlers.HandleStripeWebhook(cfg))

	r.Route("/", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware(cfg.JWTSecret))
		r.Get("/cart", handlers.GetCart)
		r.Post("/cart/items", handlers.AddToCart)
		r.Put("/cart/items/{id}", handlers.UpdateCartItem)
		r.Delete("/cart/items/{id}", handlers.RemoveFromCart)
		r.Delete("/cart", handlers.ClearCart)

		r.Post("/checkout", handlers.CreateCheckout(cfg))
		r.Get("/orders/{id}", handlers.GetOrder)
	})

	log.Println("Listening on port 8080")
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)

	if err != nil {
		log.Fatal("Failed listening on port 8080: ", err)
		return
	}
}
