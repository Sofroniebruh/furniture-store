package main

import (
	"log"

	"cart-service/config"
	"cart-service/db"
	"cart-service/handlers"
	"cart-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db.Init(cfg)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/webhook", handlers.HandleStripeWebhook(cfg))

	api := router.Group("/")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/cart", handlers.GetCart)
		api.POST("/cart/items", handlers.AddToCart)
		api.PUT("/cart/items/:id", handlers.UpdateCartItem)
		api.DELETE("/cart/items/:id", handlers.RemoveFromCart)
		api.DELETE("/cart", handlers.ClearCart)

		api.POST("/checkout", handlers.CreateCheckout(cfg))
		api.GET("/orders/:id", handlers.GetOrder)
	}

	log.Printf("Cart service starting on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}