package interceptor

import (
	"encoding/json"
	"net/http"

	"github.com/prajwalbharadwajbm/broker/internal/dtos"
)

func SendErrorResponse(w http.ResponseWriter, errorCode string, statusCode int) {
	errorMessage := errors[errorCode]
	response := dtos.InterceptorResponse{
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := dtos.InterceptorResponse{
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
