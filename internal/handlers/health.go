package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prajwalbharadwajbm/broker/internal/config"
	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
)

// Health is required to check
//   - to ensure reliability and proper functioning
//   - ensuring application availability: this is useful if hosted on cloud
//     as azure function app or aws ec2 for this instance to monitor availability.
func Health(w http.ResponseWriter, r *http.Request) {
	dbClient := db.GetProtectedClient()
	circuitBreakerState := dbClient.GetCircuitBreakerState()
	failureCount := dbClient.GetFailures()

	// Test database connectivity through circuit breaker
	dbStatus := "healthy"
	dbError := dbClient.Ping()
	if dbError != nil {
		dbStatus = "unhealthy"
		logger.Log.Error("Health check: database ping failed", dbError)
	}

	envelope := map[string]interface{}{
		"status": "available",
		"application-details": map[string]interface{}{
			"version":     "1.0.0",
			"environment": config.AppConfigInstance.GeneralConfig.Env,
		},
		"dependencies": map[string]interface{}{
			"database": map[string]interface{}{
				"status":                   dbStatus,
				"circuit_breaker_state":    circuitBreakerState,
				"circuit_breaker_failures": failureCount,
			},
		},
	}

	healthObj, err := json.Marshal(envelope)
	if err != nil {
		logger.Log.Error("failed to marshal health check response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set appropriate status code based on dependencies
	statusCode := http.StatusOK
	if dbStatus == "unhealthy" || circuitBreakerState == "OPEN" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(healthObj)
}
