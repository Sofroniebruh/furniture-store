package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
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
	log.Println(userId)
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
