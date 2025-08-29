package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"user-related-service/config"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) { return config.JWT_SECRET, nil })
}

func RetrieveIdFromCookie(r *http.Request, cookieName string) (uuid.UUID, error) {
	cookie, err := r.Cookie(cookieName)

	if err != nil || cookie == nil || cookie.Value == "" {
		return uuid.Nil, errors.New("invalid cookie")
	}

	token, err := ParseToken(cookie.Value)

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims := token.Claims.(*Claims)
	id := claims.UserID
	return id, nil
}
