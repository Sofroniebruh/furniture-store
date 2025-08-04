package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID                        uuid.UUID `json:"id" db:"id"`
	Username                  string    `json:"username" db:"username"`
	Email                     string    `json:"email" db:"email"`
	Password                  string    `json:"-" db:"password"`
	RefreshToken              *string   `json:"-" db:"refresh_token"`
	CreatedAt                 time.Time `json:"-" db:"created_at"`
	VerificationOnSignupToken string    `json:"-" db:"verification_on_signup_token"`
}
