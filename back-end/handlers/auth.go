package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"furniture-store-backend/config"
	"furniture-store-backend/db"
	"furniture-store-backend/models"
	"furniture-store-backend/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userInfoRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	var existingUser models.User
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
		})
		return
	}

	if userInfoRequest.Username == "" || userInfoRequest.Email == "" || userInfoRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Username, Email and Password are required",
		})
		return
	}

	err = db.DB.Get(&existingUser, "SELECT id FROM users WHERE email = $1", userInfoRequest.Email)

	if existingUser.ID != uuid.Nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "User already exists",
		})
		return
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error checking existing user",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfoRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	err = db.DB.QueryRow(`
		INSERT INTO users (username, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`, userInfoRequest.Username, userInfoRequest.Email, string(hashedPassword), time.Now().UTC(),
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to save user",
		})
		return
	}

	accessToken, _ := utils.GenerateToken(user.ID, config.ACCESS_TOKEN_TTL)
	refreshToken, _ := utils.GenerateToken(user.ID, config.REFRESH_TOKEN_TTL)

	_, err = db.DB.Exec("UPDATE users SET refresh_token = $1 WHERE id = $2", refreshToken, user.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to save refresh token",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   int(config.ACCESS_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(config.REFRESH_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})

	user.Username = userInfoRequest.Username
	user.Email = userInfoRequest.Email

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"created": user,
	})
}
