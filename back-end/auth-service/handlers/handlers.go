package handlers

import (
	"auth-service/config"
	"auth-service/db"
	"auth-service/models"
	"auth-service/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
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
var rabbitData struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var existingUser models.User
	var user models.User
	conn, ch, err := config.InitRabbitMq()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userInfoRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read request",
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

	verificationEmailToken, _ := utils.GenerateToken(user.ID, config.ACCESS_TOKEN_TTL)
	rabbitData.Email = userInfoRequest.Email
	rabbitData.Token = verificationEmailToken

	err = config.ProduceMessage(ch, "verifyEmail", rabbitData)

	defer conn.Close()
	defer ch.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to send email: " + err.Error(),
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

	_, err = db.DB.Exec("UPDATE users SET refresh_token = $1, verification_on_signup_token = $2 WHERE id = $3", refreshToken, verificationEmailToken, user.ID)

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
	user.VerificationOnSignupToken = verificationEmailToken

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"created": user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&userInfoRequest)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read request",
		})
	}

	if userInfoRequest.Username == "" || userInfoRequest.Email == "" || userInfoRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Username, Email and Password are required",
		})
		return
	}

	err = db.DB.Get(&user, "SELECT * FROM users WHERE email = $1", userInfoRequest.Email)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfoRequest.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	if user.VerificationOnSignupToken != "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Your email is not verified",
		})
		return
	}

	accessToken, _ := utils.GenerateToken(user.ID, config.ACCESS_TOKEN_TTL)
	refreshToken, _ := utils.GenerateToken(user.ID, config.REFRESH_TOKEN_TTL)

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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"user": user,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.RetrieveIdFromCookie(r, "refresh_token")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	_, err = db.DB.Exec("UPDATE users SET refresh_token = NULL WHERE id = $1", userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update the user",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.RetrieveIdFromCookie(r, "refresh_token")

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	refreshToken, _ := utils.GenerateToken(userId, config.REFRESH_TOKEN_TTL)
	accessToken, _ := utils.GenerateToken(userId, config.ACCESS_TOKEN_TTL)

	_, err = db.DB.Exec("UPDATE users SET refresh_token = $1 WHERE id = $2", refreshToken, userId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update the user",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(config.REFRESH_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   int(config.ACCESS_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Refreshed successfully",
	})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user := models.User{}

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Token is required",
		})
		return
	}

	parsedToken, err := utils.ParseToken(token)

	if err != nil || !parsedToken.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	id, _ := uuid.Parse(claims["sub"].(string))

	err = db.DB.Get(&user, "SELECT * from users WHERE id = $1", id)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to find user",
		})
		return
	}

	_, err = db.DB.Exec("UPDATE users SET verification_on_signup_token = NULL WHERE id = $1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update the user",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"updated": user,
	})
}
