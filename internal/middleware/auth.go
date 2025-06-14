package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/service/auth"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Log.Info("missing authorization header")
			interceptor.SendErrorResponse(w, "BPB010", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Info("invalid authorization header format")
			interceptor.SendErrorResponse(w, "BPB013", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				logger.Log.Info("token expired: %v", err)
				interceptor.SendErrorResponse(w, "BPB012", http.StatusUnauthorized)
				return
			}

			logger.Log.Info("invalid token: %v", err)
			interceptor.SendErrorResponse(w, "BPB012", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
