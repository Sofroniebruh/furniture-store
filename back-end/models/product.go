package models

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Amount      int       `json:"amount" db:"amount"`
	Price       float64   `json:"price" db:"price"`
	PictureUrl  string    `json:"pictureUrl" db:"picture_url"`
}
