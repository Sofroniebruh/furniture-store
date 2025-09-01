package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Colors      []Color        `json:"colors"`
}

type Color struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}
