package utils

import (
	"errors"
	"furniture-store-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func GenerateToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWT_SECRET)
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})
}

func RetrieveIdFromCookie(r *http.Request, cookieName string) (uuid.UUID, error) {
	cookie, err := r.Cookie(cookieName)

	if err != nil || cookie == nil || cookie.Value == "" {
		return uuid.Nil, errors.New("invalid cookie")
	}

	token, _ := ParseToken(cookie.Value)
	claims := token.Claims.(jwt.MapClaims)
	return uuid.Parse(claims["sub"].(string))
}
