package handlers

import (
	"net/http"

	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/service/auth"
	"github.com/prajwalbharadwajbm/broker/internal/utils"
)

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken handles the refresh token endpoint
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestData, err := utils.FetchDataFromRequestBody[RefreshTokenRequest](r)
	if err != nil {
		logger.Log.Error("unable to fetch request body", err)
		interceptor.SendErrorResponse(w, "BPB001", http.StatusBadRequest)
		return
	}

	if requestData.RefreshToken == "" {
		logger.Log.Info("refresh token is required")
		interceptor.SendErrorResponse(w, "BPB010", http.StatusBadRequest)
		return
	}

	// Validate refresh token exists and is not expired
	storedToken, err := repository.ValidateRefreshToken(ctx, requestData.RefreshToken)
	if err != nil {
		logger.Log.Error("failed to validate refresh token", err)
		interceptor.SendErrorResponse(w, "BPB011", http.StatusInternalServerError)
		return
	}

	if storedToken == nil {
		logger.Log.Info("invalid or expired refresh token")
		interceptor.SendErrorResponse(w, "BPB012", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	newAccessToken, err := auth.GenerateToken(storedToken.UserID.String())
	if err != nil {
		logger.Log.Error("failed to generate new access token", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Generate new refresh token for rotation
	newRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		logger.Log.Error("failed to generate new refresh token", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Update the refresh token in database (token rotation)
	err = repository.RevokeRefreshToken(ctx, requestData.RefreshToken)
	if err != nil {
		logger.Log.Error("failed to revoke old refresh token", err)
		interceptor.SendErrorResponse(w, "BPB011", http.StatusInternalServerError)
		return
	}

	refreshTokenExpiry := auth.GetRefreshTokenExpiration()
	_, err = repository.CreateRefreshToken(ctx, storedToken.UserID, newRefreshToken, refreshTokenExpiry)
	if err != nil {
		logger.Log.Error("failed to store new refresh token", err)
		interceptor.SendErrorResponse(w, "BPB011", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // access token expires in 15 minutes in seconds
	}

	logger.Log.Infof("Successfully refreshed token for user_id: %s", storedToken.UserID.String())
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}

// RevokeRefreshToken handles token revocation (logout)
func RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestData, err := utils.FetchDataFromRequestBody[RefreshTokenRequest](r)
	if err != nil {
		logger.Log.Error("unable to fetch request body", err)
		interceptor.SendErrorResponse(w, "BPB001", http.StatusBadRequest)
		return
	}

	if requestData.RefreshToken == "" {
		logger.Log.Info("refresh token is required")
		interceptor.SendErrorResponse(w, "BPB010", http.StatusBadRequest)
		return
	}

	err = repository.RevokeRefreshToken(ctx, requestData.RefreshToken)
	if err != nil {
		logger.Log.Error("failed to revoke refresh token", err)
		interceptor.SendErrorResponse(w, "BPB011", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Token revoked successfully",
	}

	logger.Log.Info("Successfully revoked refresh token")
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}
