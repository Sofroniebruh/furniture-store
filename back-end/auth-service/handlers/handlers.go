package handlers

import (
	"auth-service/config"
	"auth-service/db"
	"auth-service/models"
	"auth-service/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type userInfoRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userCodeRequest struct {
	Code  string `json:"code"`
	Email string `json:"email"`
}

var rabbitData struct {
	Email       string `json:"email"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var existingUser models.User
	var user models.User
	var userInfo userInfoRequest

	err := json.NewDecoder(r.Body).Decode(&userInfo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read request",
		})
		return
	}

	if userInfo.Email == "" || userInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Email and Password are required",
		})
		return
	}

	err = db.DB.Get(&existingUser, "SELECT id FROM users WHERE email = $1", userInfo.Email)

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

	data, err := json.Marshal(userInfo)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error marshalling user info",
		})
		return
	}

	go func() {
		response, err := GenerateCode(data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Error generating code: " + err.Error(),
			})
			return
		}

		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Error sending email",
			})
			return
		}
	}()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	err = db.DB.QueryRow(`
		INSERT INTO users (email, password, created_at, is_email_verified)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`, userInfo.Email, string(hashedPassword), time.Now().UTC(), false,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to save user",
		})
		return
	}

	user.Email = userInfo.Email
	user.IsEmailVerified = false

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"created": user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var userInfo userInfoRequest

	err := json.NewDecoder(r.Body).Decode(&userInfo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read request",
		})
	}

	if userInfo.Email == "" || userInfo.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Email and Password are required",
		})
		return
	}

	err = db.DB.Get(&user, "SELECT * FROM users WHERE email = $1", userInfo.Email)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Unauthorized",
		})
		return
	}

	if !user.IsEmailVerified {
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

func Verify(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var userWithCode userCodeRequest

	err := json.NewDecoder(r.Body).Decode(&userWithCode)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
		})
		return
	}

	var rdb = config.NewRedisConfig()
	code, err := rdb.Get(userWithCode.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get the code",
		})
		return
	}

	if code != userWithCode.Code {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid code",
		})
		return
	}

	err = rdb.Delete(userWithCode.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete the code",
		})
		return
	}

	err = db.DB.QueryRow("UPDATE users SET is_email_verified = $1 WHERE email = $2 RETURNING id", true, userWithCode.Email).Scan(&user.ID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to update the user",
		})
		return
	}

	user.Email = userWithCode.Email
	user.IsEmailVerified = true

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

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"verified": user,
	})
}

func ResendCode(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}
	var user userEmail
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request",
		})
		return
	}

	data, _ := json.Marshal(user)

	response, err := GenerateCode(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to generate code",
		})
		return
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Error sending email",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully sent email",
	})
}

func GenerateCode(body []byte) (*config.ResponsePythonHandler, error) {
	var user userInfoRequest
	err := json.Unmarshal(body, &user)
	minRange := 10000
	maxRange := 99999

	if err != nil {
		return &config.ResponsePythonHandler{}, errors.New("Failed to unmarshal JSON: " + err.Error())
	}

	rdb := config.NewRedisConfig()
	randomCode := rand.Intn(maxRange-minRange+1) + minRange
	stringRandomNumber := strconv.Itoa(int(randomCode))

	err = rdb.Set(user.Email, stringRandomNumber, 120*time.Second)

	if err != nil {
		return &config.ResponsePythonHandler{}, errors.New("Failed to set random code: " + err.Error())
	}

	response, err := SendEmail(stringRandomNumber, user.Email)

	if err != nil {
		return &config.ResponsePythonHandler{}, errors.New("Failed to send email: " + err.Error())
	}

	return &response, nil
}

func SendEmail(data string, email string) (config.ResponsePythonHandler, error) {
	conn, ch, err := config.InitRabbitMq()

	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	requestQueue, _ := config.DeclareQueue(
		ch,
		"verifyEmail",
		false,
		false,
		false,
		false)

	responseQueue, _ := config.DeclareQueue(
		ch,
		"responseFromRequestQueue",
		false,
		false,
		false,
		false)

	rabbitData.Email = email
	rabbitData.MessageBody = data
	rabbitData.Subject = "Verify your email"

	correlationID := uuid.New()

	err = config.ProduceMessage(requestQueue, responseQueue, ch, rabbitData, correlationID)

	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	response, err := config.WaitForResponseQueue(ch, responseQueue.Name, correlationID)

	conn.Close()
	ch.Close()

	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return config.ResponsePythonHandler{}, errors.New("failed to send the email")
	}

	return *response, nil
}
