package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"cart-service/config"
	"cart-service/models"
	"cart-service/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateCheckout(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

		var req models.CheckoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		order, err := services.CreatePaymentIntent(userID, cfg)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			if err.Error() == "cart is empty" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			if err.Error() == "insufficient stock" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create checkout session"})
			return
		}

		response := models.CheckoutResponse{
			OrderID:      order.ID,
			ClientSecret: order.StripeClientSecret,
			Status:       "created",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func HandleStripeWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unable to read request body"})
			return
		}

		signature := r.Header.Get("Stripe-Signature")
		if signature == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Missing Stripe-Signature header"})
			return
		}

		err = services.HandleStripeWebhook(body, signature, cfg)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Webhook processed successfully"})
	}
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	orderIDStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid order ID"})
		return
	}

	order, err := services.GetOrder(userID, orderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}