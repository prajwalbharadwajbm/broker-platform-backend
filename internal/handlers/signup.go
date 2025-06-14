package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/dtos"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/utils"
	"github.com/prajwalbharadwajbm/broker/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// register a new user with email and password.
func Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userData, err := utils.FetchDataFromRequestBody[dtos.User](r)
	if err != nil {
		logger.Log.Error("unable to fetch request body", err)
		interceptor.SendErrorResponse(w, "BPB001", http.StatusBadRequest)
		return
	}

	valid, err := validateRequestBody(userData)
	if !valid || err != nil {
		logger.Log.Info("request body is not valid %v", err)
		interceptor.SendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := registerUser(ctx, userData)
	if err != nil {
		logger.Log.Error("unable to register user", err)
		interceptor.SendErrorResponse(w, "BPB005", http.StatusBadRequest)
		return
	}
	response := map[string]interface{}{
		"userID": userID,
	}
	logger.Log.Infof("Successfully registered user_id: %s", userID)
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}

func registerUser(ctx context.Context, userData dtos.User) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash password: %w", err)
	}
	userId, err := repository.AddUser(ctx, userData.Email, hashedPassword)
	if err != nil {
		return "", fmt.Errorf("unable to add user: %w", err)
	}
	return userId, nil
}

func validateRequestBody(userData dtos.User) (bool, error) {
	if valid, err := validator.IsValidEmail(userData.Email); !valid || err != nil {
		return false, fmt.Errorf("invalid email: %w", err)
	}

	if valid, err := validator.IsValidPassword(userData.Email, userData.Password); !valid || err != nil {
		return false, fmt.Errorf("invalid password: %w", err)
	}
	return true, nil
}
