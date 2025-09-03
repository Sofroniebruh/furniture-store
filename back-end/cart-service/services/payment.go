package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"cart-service/config"
	"cart-service/db"
	"cart-service/models"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

func CreatePaymentIntent(userID uuid.UUID, cfg *config.Config) (*models.Order, error) {
	stripe.Key = cfg.StripeSecretKey

	total, err := GetCartTotal(userID)
	if err != nil {
		return nil, err
	}

	if total <= 0 {
		return nil, errors.New("cart is empty")
	}

	err = ValidateCartForCheckout(userID)
	if err != nil {
		return nil, err
	}

	orderID := uuid.New()
	_, err = db.DB.Exec(
		"INSERT INTO orders (id, user_id, total_amount, status) VALUES ($1, $2, $3, $4)",
		orderID, userID, total, models.OrderStatusPending,
	)
	if err != nil {
		return nil, err
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(total * 100)), // Convert to cents
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Metadata: map[string]string{
			"order_id": orderID.String(),
			"user_id":  userID.String(),
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		db.DB.Exec("DELETE FROM orders WHERE id = $1", orderID)
		return nil, fmt.Errorf("failed to create payment intent: %v", err)
	}

	_, err = db.DB.Exec(
		"UPDATE orders SET stripe_payment_id = $1 WHERE id = $2",
		pi.ID, orderID,
	)
	if err != nil {
		log.Printf("Failed to update order with payment intent ID: %v", err)
	}

	err = createOrderItemsFromCart(userID, orderID)
	if err != nil {
		log.Printf("Failed to create order items: %v", err)
	}

	return &models.Order{
		ID:                 orderID,
		UserID:             userID,
		StripePaymentID:    pi.ID,
		TotalAmount:        total,
		Status:             models.OrderStatusPending,
		StripeClientSecret: pi.ClientSecret,
	}, nil
}

func createOrderItemsFromCart(userID, orderID uuid.UUID) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		SELECT $1, ci.product_id, ci.quantity, p.price
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = $2
	`
	
	_, err := db.DB.Exec(query, orderID, userID)
	return err
}

func HandleStripeWebhook(body []byte, signature string, cfg *config.Config) error {
	event, err := webhook.ConstructEvent(body, signature, cfg.StripeWebhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %v", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &pi)
		if err != nil {
			return fmt.Errorf("error parsing webhook JSON: %v", err)
		}

		err = handlePaymentSuccess(&pi, cfg)
		if err != nil {
			log.Printf("Error handling payment success: %v", err)
			return err
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &pi)
		if err != nil {
			return fmt.Errorf("error parsing webhook JSON: %v", err)
		}

		err = handlePaymentFailure(&pi)
		if err != nil {
			log.Printf("Error handling payment failure: %v", err)
			return err
		}

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

func handlePaymentSuccess(pi *stripe.PaymentIntent, cfg *config.Config) error {
	orderID := pi.Metadata["order_id"]
	userIDStr := pi.Metadata["user_id"]

	if orderID == "" || userIDStr == "" {
		return errors.New("missing order metadata in payment intent")
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %v", err)
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	_, err = db.DB.Exec(
		"UPDATE orders SET status = $1 WHERE id = $2",
		models.OrderStatusPaid, orderUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	err = reduceProductStock(orderUUID)
	if err != nil {
		log.Printf("Failed to reduce product stock: %v", err)
	}

	err = addToHistory(userUUID, orderUUID, cfg)
	if err != nil {
		log.Printf("Failed to add items to history: %v", err)
	}

	err = ClearCart(userUUID)
	if err != nil {
		log.Printf("Failed to clear cart: %v", err)
	}

	log.Printf("Successfully processed payment for order %s", orderID)
	return nil
}

func handlePaymentFailure(pi *stripe.PaymentIntent) error {
	orderID := pi.Metadata["order_id"]
	if orderID == "" {
		return errors.New("missing order metadata in payment intent")
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID: %v", err)
	}

	_, err = db.DB.Exec(
		"UPDATE orders SET status = $1 WHERE id = $2",
		models.OrderStatusCancelled, orderUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	log.Printf("Payment failed for order %s", orderID)
	return nil
}

func reduceProductStock(orderID uuid.UUID) error {
	query := `
		UPDATE products 
		SET stock = stock - oi.quantity
		FROM order_items oi
		WHERE products.id = oi.product_id AND oi.order_id = $1
	`
	
	_, err := db.DB.Exec(query, orderID)
	return err
}

func addToHistory(userID, orderID uuid.UUID, cfg *config.Config) error {
	query := `SELECT product_id FROM order_items WHERE order_id = $1`
	rows, err := db.DB.Query(query, orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID uuid.UUID
		if err := rows.Scan(&productID); err != nil {
			log.Printf("Error scanning product ID: %v", err)
			continue
		}

		err = callAddToHistory(userID, productID, cfg)
		if err != nil {
			log.Printf("Failed to add product %s to history for user %s: %v", productID, userID, err)
		}
	}

	return nil
}

func callAddToHistory(userID, productID uuid.UUID, cfg *config.Config) error {
	url := fmt.Sprintf("%s/internal/user/history", cfg.UserRelatedServiceURL)
	
	payload := map[string]string{
		"user_id":    userID.String(),
		"product_id": productID.String(),
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to add to history: status code %d", resp.StatusCode)
	}
	
	return nil
}

func GetOrder(userID, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := db.DB.Get(&order,
		"SELECT * FROM orders WHERE id = $1 AND user_id = $2",
		orderID, userID,
	)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
			   p.name, p.description, p.stock, p.event, p.model
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`
	
	rows, err := db.DB.Query(query, orderID)
	if err != nil {
		return &order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		var product models.Product
		
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price,
			&product.Name, &product.Description, &product.Stock, &product.Event, &product.Model,
		)
		if err != nil {
			return &order, err
		}
		
		product.ID = item.ProductID
		product.Price = item.Price
		item.Product = &product
		order.Items = append(order.Items, item)
	}

	return &order, nil
}