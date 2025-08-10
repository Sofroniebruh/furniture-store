package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"products-service/config"
	"products-service/utils"
)

func ProtectedWithRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, roles, err := utils.RetrieveIdAndRoleFromCookie(r, "access_token")

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Unauthorized",
				})
				return
			}

			hasRole := false

			for _, role := range roles {
				for _, allowedRole := range allowedRoles {
					if role == allowedRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
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
}
