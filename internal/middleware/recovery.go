package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
)

// RecoveryMiddleware is a global panic recovery middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error(
					fmt.Sprintf("Panic recovered: %v", err),
					fmt.Errorf("panic: %v", err),
				)

				// Log the stack trace for debugging
				logger.Log.Error("Stack trace", fmt.Errorf("%s", debug.Stack()))

				// Send error response to client
				interceptor.SendErrorResponse(w, "BPB500", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
