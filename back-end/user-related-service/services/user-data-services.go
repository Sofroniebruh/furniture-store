package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"user-related-service/config"
	"user-related-service/db"
	"user-related-service/models"
)

type UserDataWithCodeResponse struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type UserDataWithEmail struct {
	Email string `json:"email"`
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(config.UserIdKey).(uuid.UUID)
	var user models.User

	err := db.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userId)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]models.User{
		"user": user,
	})
}

func StoreCode(w http.ResponseWriter, r *http.Request) {
	var requestBody UserDataWithEmail
	conn, ch, err := config.InitRabbitMq()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to connect to RabbitMQ",
		})
		return
	}

	queue, _ := config.DeclareQueue(ch, "codeQueue", false, false, false, false)
	replyQueue, _ := config.DeclareQueue(ch, "replyCodeQueue", false, false, false, false)

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	minRange := 100000
	maxRange := 999999

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	rdb := config.NewRedisConfig()
	randomCode := rand.Intn(maxRange-minRange+1) + minRange
	stringRandomNumber := strconv.Itoa(int(randomCode))

	err = rdb.Set(requestBody.Email, stringRandomNumber)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to save the code",
		})
		return
	}

	correlationID := uuid.New()

	var body struct {
		Code  string `json:"code"`
		Email string `json:"email"`
	}

	body.Code = stringRandomNumber
	body.Email = requestBody.Email

	log.Println("Body: ", body)

	err = config.ProduceMessage(queue, replyQueue, ch, body, correlationID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to send message with RabbitMQ",
		})
		return
	}

	response, err := config.WaitForResponseQueue(ch, "replyCodeQueue", correlationID)

	conn.Close()
	ch.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to send email: " + err.Error(),
		})
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Received nil response from queue",
		})
		return
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to send email: " + response.Message,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"email": requestBody.Email,
		"code":  randomCode,
	})
}

func CompareCode(w http.ResponseWriter, r *http.Request) {
	var requestBody UserDataWithCodeResponse
	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	rdb := config.NewRedisConfig()
	code, err := rdb.Get(requestBody.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get the code",
		})
		return
	}

	if code != requestBody.Code {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Code does not match",
		})
		return
	}

	err = rdb.Delete(requestBody.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to delete the code",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "User code is correct",
	})
}
