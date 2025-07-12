package middleware

import (
	"context"
	"encoding/json"
	"furniture-store-backend/config"
	"furniture-store-backend/utils"
	"net/http"
)

func Protected(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.RetrieveIdFromCookie(r, "access_token")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		ctx := context.WithValue(r.Context(), config.UserIdKey, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
