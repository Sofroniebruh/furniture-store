package models

import (
	"time"

	"github.com/google/uuid"
)

type StockHistory struct {
	ID            uuid.UUID         `json:"id" db:"id"`
	ProductID     uuid.UUID         `json:"product_id" db:"product_id"`
	Type          StockHistoryType  `json:"type" db:"type"`
	Quantity      int               `json:"quantity" db:"quantity"`
	PreviousStock int               `json:"previous_stock" db:"previous_stock"`
	NewStock      int               `json:"new_stock" db:"new_stock"`
	Reason        string            `json:"reason" db:"reason"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	Product       *Product          `json:"product,omitempty"`
}

type StockHistoryType string

const (
	StockHistoryTypeIn         StockHistoryType = "in"
	StockHistoryTypeOut        StockHistoryType = "out"
	StockHistoryTypeAdjustment StockHistoryType = "adjustment"
)

type StockHistoryRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
	Reason    string    `json:"reason" binding:"required"`
}