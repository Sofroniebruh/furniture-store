package handlers

import (
	"auth-service/config"
	"auth-service/db"
	"auth-service/models"
	"auth-service/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"os"
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

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

var rabbitData struct {
	Email       string `json:"email"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
}

// Helper functions
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func respondWithSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func validateSignUpRequest(w http.ResponseWriter, request userInfoRequest) bool {
	if request.Email == "" || request.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return false
	}
	return true
}

func validateLoginRequest(w http.ResponseWriter, request userInfoRequest) bool {
	if request.Email == "" || request.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password are required")
		return false
	}
	return true
}

func validateCodeRequest(w http.ResponseWriter, request userCodeRequest) bool {
	if request.Code == "" || request.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Code and email are required")
		return false
	}
	return true
}

func validateEmailRequest(w http.ResponseWriter, email string) bool {
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return false
	}
	return true
}

func userExists(email string) bool {
	var count int
	if err := db.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", email); err != nil {
		return true
	}
	return count > 0
}

func hashPassword(w http.ResponseWriter, password string) (string, bool) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return "", false
	}
	return string(hashedPassword), true
}

func verifyPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func parseRequestBody(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return false
	}
	return true
}

func getUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	return user, err
}

func getUserById(id uuid.UUID) (models.User, error) {
	var user models.User
	err := db.DB.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

func generateAndSetTokens(w http.ResponseWriter, user models.User) bool {
	accessToken, err := utils.GenerateToken(user, config.ACCESS_TOKEN_TTL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return false
	}

	refreshToken, err := utils.GenerateToken(user, config.REFRESH_TOKEN_TTL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return false
	}

	if !updateUserRefreshToken(w, user.ID, refreshToken) {
		return false
	}

	setAuthCookies(w, accessToken, refreshToken)
	return true
}

func updateUserRefreshToken(w http.ResponseWriter, userId uuid.UUID, refreshToken string) bool {
	_, err := db.DB.Exec("UPDATE users SET refresh_token = $1 WHERE id = $2", refreshToken, userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user refresh token")
		return false
	}
	return true
}

func setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	sameSiteMode := http.SameSiteLaxMode
	secure := false

	if os.Getenv("ENVIRONMENT") == "production" {
		sameSiteMode = http.SameSiteNoneMode
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   int(config.ACCESS_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSiteMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(config.REFRESH_TOKEN_TTL.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSiteMode,
	})
}

func clearAuthCookies(w http.ResponseWriter) {
	sameSiteMode := http.SameSiteLaxMode
	secure := false

	if os.Getenv("ENVIRONMENT") == "production" {
		sameSiteMode = http.SameSiteNoneMode
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSiteMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSiteMode,
	})
}

func createUser(w http.ResponseWriter, email, hashedPassword string, roles []string, isVerified bool) (models.User, bool) {
	var user models.User
	err := db.DB.QueryRow(`
		INSERT INTO users (email, password, created_at, is_email_verified, roles)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`, email, hashedPassword, time.Now().UTC(), isVerified, pq.Array(roles),
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return user, false
	}

	user.Email = email
	user.IsEmailVerified = isVerified
	user.Roles = roles
	return user, true
}

func verifyUserEmail(w http.ResponseWriter, email string) (models.User, bool) {
	var user models.User
	err := db.DB.QueryRow("UPDATE users SET is_email_verified = $1 WHERE email = $2 RETURNING id, roles", true, email).Scan(&user.ID, &user.Roles)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update the user")
		return user, false
	}
	user.Email = email
	user.IsEmailVerified = true
	return user, true
}

func isAdminCredentials(email, password string) bool {
	return email == os.Getenv("ADMIN_EMAIL") && password == os.Getenv("ADMIN_PASSWORD")
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var userInfo userInfoRequest

	if !parseRequestBody(w, r, &userInfo) {
		return
	}

	if !validateSignUpRequest(w, userInfo) {
		return
	}

	if isAdminCredentials(userInfo.Email, userInfo.Password) {
		HandleAdminSignUp(userInfo, w)
		return
	}

	if userExists(userInfo.Email) {
		respondWithError(w, http.StatusBadRequest, "User already exists")
		return
	}

	data, err := json.Marshal(userInfo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error marshalling user info")
		return
	}

	go func() {
		response, err := GenerateCode(data)
		if err != nil {
			log.Printf("Error generating code: %v", err)
			return
		}
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			log.Printf("Error sending email: status code %d", response.StatusCode)
			return
		}
	}()

	hashedPassword, ok := hashPassword(w, userInfo.Password)
	if !ok {
		return
	}

	userRoles := []string{"user"}
	user, ok := createUser(w, userInfo.Email, hashedPassword, userRoles, false)
	if !ok {
		return
	}

	respondWithSuccess(w, http.StatusCreated, map[string]models.User{
		"created": user,
	})
}

func HandleAdminSignUp(userInfo userInfoRequest, w http.ResponseWriter) {
	if userExists(userInfo.Email) {
		respondWithError(w, http.StatusBadRequest, "User already exists")
		return
	}

	hashedPassword, ok := hashPassword(w, userInfo.Password)
	if !ok {
		return
	}

	userRoles := []string{"admin", "user"}
	user, ok := createUser(w, userInfo.Email, hashedPassword, userRoles, true)
	if !ok {
		return
	}

	if !generateAndSetTokens(w, user) {
		return
	}

	respondWithSuccess(w, http.StatusCreated, map[string]models.User{
		"user": user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var userInfo userInfoRequest

	if !parseRequestBody(w, r, &userInfo) {
		return
	}

	if !validateLoginRequest(w, userInfo) {
		return
	}

	user, err := getUserByEmail(userInfo.Email)
	if err != nil || !verifyPassword(user.Password, userInfo.Password) {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !user.IsEmailVerified {
		respondWithError(w, http.StatusUnauthorized, "Your email is not verified")
		return
	}

	if !generateAndSetTokens(w, user) {
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]models.User{
		"user": user,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.RetrieveIdFromCookie(r, "refresh_token")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = db.DB.Exec("UPDATE users SET refresh_token = NULL WHERE id = $1", userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update the user")
		return
	}

	clearAuthCookies(w)

	respondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.RetrieveIdFromCookie(r, "refresh_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !generateAndSetTokens(w, user) {
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Refreshed successfully",
	})
}

func Verify(w http.ResponseWriter, r *http.Request) {
	var userWithCode userCodeRequest

	if !parseRequestBody(w, r, &userWithCode) {
		return
	}

	if !validateCodeRequest(w, userWithCode) {
		return
	}

	rdb := config.NewRedisConfig()
	code, err := rdb.Get(userWithCode.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get the code")
		return
	}

	if code != userWithCode.Code {
		respondWithError(w, http.StatusUnauthorized, "Invalid code")
		return
	}

	err = rdb.Delete(userWithCode.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete the code")
		return
	}

	log.Println(userWithCode.Email)
	user, ok := verifyUserEmail(w, userWithCode.Email)
	if !ok {
		return
	}

	if !generateAndSetTokens(w, user) {
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]models.User{
		"verified": user,
	})
}

func ResendCode(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}
	var user userEmail

	if !parseRequestBody(w, r, &user) {
		return
	}

	if !validateEmailRequest(w, user.Email) {
		return
	}

	data, _ := json.Marshal(user)

	response, err := GenerateCode(data)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate code")
		return
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		respondWithError(w, http.StatusInternalServerError, "Error sending email")
		return
	}

	respondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Successfully sent email",
	})
}

func GenerateCode(body []byte) (*config.ResponsePythonHandler, error) {
	var user userInfoRequest
	err := json.Unmarshal(body, &user)
	if err != nil {
		return &config.ResponsePythonHandler{}, errors.New("Failed to unmarshal JSON: " + err.Error())
	}

	rdb := config.NewRedisConfig()
	randomCode := rand.Intn(90000) + 10000
	stringRandomNumber := strconv.Itoa(randomCode)

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
	defer conn.Close()
	defer ch.Close()

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
	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return config.ResponsePythonHandler{}, errors.New("failed to send the email")
	}

	return *response, nil
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request forgotPasswordRequest

	if !parseRequestBody(w, r, &request) {
		return
	}

	if !validateEmailRequest(w, request.Email) {
		return
	}

	if !userExists(request.Email) {
		respondWithError(w, http.StatusNotFound, "No account found with this email address")
		return
	}

	resetToken := uuid.New().String()

	rdb := config.NewRedisConfig()
	err := rdb.Set("reset_"+request.Email, resetToken, 30*time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate reset token")
		return
	}

	go func() {
		response, err := SendPasswordResetEmail(resetToken, request.Email)
		if err != nil {
			log.Printf("Error sending password reset email: %v", err)
			return
		}
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			log.Printf("Error sending password reset email: status code %d", response.StatusCode)
			return
		}
	}()

	respondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Password reset email sent successfully",
	})
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var request resetPasswordRequest

	if !parseRequestBody(w, r, &request) {
		return
	}

	if request.Token == "" || request.NewPassword == "" {
		respondWithError(w, http.StatusBadRequest, "Token and new password are required")
		return
	}

	rdb := config.NewRedisConfig()
	var userEmail string
	found := false

	keys, err := rdb.Keys("reset_*")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to verify reset token")
		return
	}

	for _, key := range keys {
		storedToken, err := rdb.Get(key)
		if err != nil {
			continue
		}
		if storedToken == request.Token {
			userEmail = key[6:]
			found = true
			break
		}
	}

	if !found {
		respondWithError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}

	hashedPassword, ok := hashPassword(w, request.NewPassword)
	if !ok {
		return
	}

	_, err = db.DB.Exec("UPDATE users SET password = $1 WHERE email = $2", hashedPassword, userEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	err = rdb.Delete("reset_" + userEmail)
	if err != nil {
		log.Printf("Failed to delete reset token: %v", err)
	}

	respondWithSuccess(w, http.StatusOK, map[string]string{
		"message": "Password reset successfully",
	})
}

func SendPasswordResetEmail(token string, email string) (config.ResponsePythonHandler, error) {
	conn, ch, err := config.InitRabbitMq()
	if err != nil {
		return config.ResponsePythonHandler{}, err
	}
	defer conn.Close()
	defer ch.Close()

	requestQueue, _ := config.DeclareQueue(
		ch,
		"resetPasswordEmail",
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

	baseURL := os.Getenv("FRONTEND_URL")
	if baseURL == "" {
		baseURL = "http://localhost:5173"
	}
	resetLink := baseURL + "/reset-password?token=" + token

	rabbitData.Email = email
	rabbitData.MessageBody = resetLink
	rabbitData.Subject = "Reset your password"

	correlationID := uuid.New()

	err = config.ProduceMessage(requestQueue, responseQueue, ch, rabbitData, correlationID)
	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	response, err := config.WaitForResponseQueue(ch, responseQueue.Name, correlationID)
	if err != nil {
		return config.ResponsePythonHandler{}, err
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		return config.ResponsePythonHandler{}, errors.New("failed to send the email")
	}

	return *response, nil
}
