package utils

import (
	"auth-service/config"
	"auth-service/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Roles  []string  `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Roles:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fumi",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWT_SECRET)
}

func ParseToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_SECRET, nil
	})
}

func RetrieveIdFromCookie(r *http.Request, cookieName string) (uuid.UUID, error) {
	cookie, err := r.Cookie(cookieName)

	if err != nil || cookie == nil || cookie.Value == "" {
		return uuid.Nil, errors.New("invalid cookie")
	}

	token, err := ParseToken(cookie.Value)

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("unauthorized")
	}

	claims := token.Claims.(*Claims)
	return uuid.Parse(claims.UserID.String())
}
