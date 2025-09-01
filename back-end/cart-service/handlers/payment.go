package handlers

import (
	"io"
	"net/http"

	"cart-service/config"
	"cart-service/models"
	"cart-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCheckout(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

		var req models.CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order, err := services.CreatePaymentIntent(userID, cfg)
		if err != nil {
			if err.Error() == "cart is empty" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err.Error() == "insufficient stock" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
			return
		}

		response := models.CheckoutResponse{
			OrderID:      order.ID,
			ClientSecret: order.StripeClientSecret,
			Status:       "created",
		}

		c.JSON(http.StatusOK, response)
	}
}

func HandleStripeWebhook(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
			return
		}

		signature := c.GetHeader("Stripe-Signature")
		if signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
			return
		}

		err = services.HandleStripeWebhook(body, signature, cfg)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
	}
}

func GetOrder(c *gin.Context) {
	userID := c.MustGet(string(config.UserIdKey)).(uuid.UUID)

	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := services.GetOrder(userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}