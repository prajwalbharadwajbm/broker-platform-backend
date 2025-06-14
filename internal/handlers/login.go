package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/dtos"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/service/auth"
	"github.com/prajwalbharadwajbm/broker/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userData, err := utils.FetchDataFromRequestBody[dtos.User](r)
	if err != nil {
		logger.Log.Error("unable to fetch request body", err)
		interceptor.SendErrorResponse(w, "BPB001", http.StatusBadRequest)
		return
	}

	authenticated, userId, err := authenticateUser(ctx, userData)
	if err != nil {
		logger.Log.Error("authentication error", err)
		interceptor.SendErrorResponse(w, "BPB006", http.StatusInternalServerError)
		return
	}
	if !authenticated {
		logger.Log.Info("invalid username or password")
		if errors.Is(err, errors.New("user not found")) {
			interceptor.SendErrorResponse(w, "BPB007", http.StatusNotFound)
		} else {
			interceptor.SendErrorResponse(w, "BPB008", http.StatusUnauthorized)
		}
		return
	}

	// Generate both access and refresh tokens
	tokenPair, err := auth.GenerateTokenPair(userId)
	if err != nil {
		logger.Log.Error("failed to generate token pair", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.Error("failed to parse user ID", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Store refresh token in database
	refreshTokenExpiry := auth.GetRefreshTokenExpiration()
	_, err = repository.CreateRefreshToken(ctx, userUUID, tokenPair.RefreshToken, refreshTokenExpiry)
	if err != nil {
		logger.Log.Error("failed to store refresh token", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // access token expires in 15 minutes in seconds
		"user_id":       userId,
	}
	logger.Log.Infof("Successfully logged in user_id: %s", userId)
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}

func authenticateUser(ctx context.Context, userData dtos.User) (bool, string, error) {
	userId, hashedPassword, err := repository.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return false, "", fmt.Errorf("unable to fetch user by email: %w", err)
	}
	if userId == "" {
		return false, "", nil
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(userData.Password))
	if err != nil {
		return false, "", nil
	}
	return true, userId, nil
}
