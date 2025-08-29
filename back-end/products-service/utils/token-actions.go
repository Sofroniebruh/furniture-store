package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"products-service/config"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})
}

func RetrieveIdAndRoleFromCookie(r *http.Request, cookieName string) (uuid.UUID, []string, error) {
	cookie, err := r.Cookie(cookieName)

	if err != nil || cookie == nil || cookie.Value == "" {
		return uuid.Nil, nil, errors.New("invalid cookie")
	}

	token, err := ParseToken(cookie.Value)

	if err != nil || !token.Valid {
		return uuid.Nil, nil, errors.New("unauthorized")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || claims == nil {
		return uuid.Nil, nil, errors.New("invalid token claims")
	}

	id := claims.UserId
	roles := claims.Roles

	return id, roles, nil
}
