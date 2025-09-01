package handlers

import (
	"net/http"

	"cart-service/config"
	"cart-service/models"
	"cart-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCart(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	cart, err := services.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func AddToCart(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddToCart(userID, req.ProductID, req.Quantity)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart successfully"})
}

func UpdateCartItem(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.UpdateCartItem(userID, itemID, req.Quantity)
	if err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "insufficient stock" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
}

func RemoveFromCart(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	err = services.RemoveFromCart(userID, itemID)
	if err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

func ClearCart(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	err := services.ClearCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
