package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"
)

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var stringArray pq.StringArray
	if err := stringArray.Scan(value); err != nil {
		return err
	}

	*s = StringSlice(stringArray)
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	return pq.StringArray(s).Value()
}

type User struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	Email           string      `json:"email" db:"email"`
	Password        string      `json:"-" db:"password"`
	Roles           StringSlice `json:"roles" db:"roles"`
	RefreshToken    *string     `json:"-" db:"refresh_token"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	IsEmailVerified bool        `json:"is_email_verified" db:"is_email_verified"`
}
