package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Product   *Product  `json:"product,omitempty"`
}

type Cart struct {
	UserID     uuid.UUID  `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPrice float64    `json:"total_price"`
}

type Order struct {
	ID                 uuid.UUID   `json:"id" db:"id"`
	UserID             uuid.UUID   `json:"user_id" db:"user_id"`
	StripePaymentID    string      `json:"stripe_payment_id" db:"stripe_payment_id"`
	TotalAmount        float64     `json:"total_amount" db:"total_amount"`
	Status             OrderStatus `json:"status" db:"status"`
	CreatedAt          time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at" db:"updated_at"`
	Items              []OrderItem `json:"items,omitempty"`
	StripeClientSecret string      `json:"stripe_client_secret,omitempty"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
	Product   *Product  `json:"product,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Product struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	Stock       int            `json:"stock" db:"stock"`
	Price       float64        `json:"price" db:"price"`
	PictureUrls pq.StringArray `json:"pictureUrls" db:"picture_urls"`
	Event       string         `json:"event" db:"event"`
	Model       string         `json:"model" db:"model"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CheckoutRequest struct {
	PaymentMethodID string `json:"payment_method_id,omitempty"`
}

type CheckoutResponse struct {
	OrderID      uuid.UUID `json:"order_id"`
	ClientSecret string    `json:"client_secret"`
	Status       string    `json:"status"`
}
