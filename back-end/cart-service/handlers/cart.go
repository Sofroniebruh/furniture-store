package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"cart-service/config"
	"cart-service/models"
	"cart-service/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	cart, err := services.GetCart(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to get cart"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	var req models.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err := services.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		log.Printf("AddToCart error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "product not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient stock" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart successfully"})
}

func UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	itemIDStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	var req models.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = services.UpdateCartItem(userID, itemID, req.Quantity)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "cart item not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient stock" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update cart item"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart item updated successfully"})
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	itemIDStr := chi.URLParam(r, "id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	err = services.RemoveFromCart(userID, itemID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "cart item not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove item from cart"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart successfully"})
}

func ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(config.UserIdKey).(uuid.UUID)

	err := services.ClearCart(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to clear cart"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared successfully"})
}
