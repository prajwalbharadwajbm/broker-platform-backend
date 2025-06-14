package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prajwalbharadwajbm/broker/internal/config"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
)

// Health is required to check
//   - to ensure reliability and proper functioning
//   - ensuring application availability: this is useful if hosted on cloud
//     as azure function app or aws ec2 for this instance to monitor availability.
func Health(w http.ResponseWriter, r *http.Request) {
	envelope := map[string]interface{}{
		"status": "available",
		"application-details": map[string]interface{}{
			"version":     "1.0.0",
			"environment": config.AppConfigInstance.GeneralConfig.Env,
		},
	}
	healthObj, err := json.Marshal(envelope)
	if err != nil {
		logger.Log.Error("failed to marshal health check response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(healthObj)
}
