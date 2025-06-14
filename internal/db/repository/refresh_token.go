package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	circuit "github.com/rubyist/circuitbreaker"
)

// CreateRefreshToken stores a new refresh token in the database
func CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (*models.RefreshToken, error) {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) 
			  VALUES ($1, $2, $3)`

	_, err := db.ExecContext(dbCtx, query, refreshToken.UserID,
		refreshToken.Token, refreshToken.ExpiresAt)
	if err != nil {
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("Create refresh token blocked by circuit breaker", err)
			return nil, errors.New("authentication service temporarily unavailable")
		}
		return nil, err
	}

	return refreshToken, nil
}

// ValidateRefreshToken retrieves and validates a refresh token by its token value
// Returns nil if token is not found or expired
func ValidateRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var refreshToken models.RefreshToken

	query := `SELECT user_id, token, expires_at 
			  FROM refresh_tokens 
			  WHERE token = $1 AND expires_at > NOW() AT TIME ZONE 'UTC'`

	row, err := db.QueryRowContext(dbCtx, query, token)
	if err != nil {
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("Validate refresh token blocked by circuit breaker", err)
			return nil, errors.New("authentication service temporarily unavailable")
		}
		return nil, err
	}

	err = row.Scan(
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Token not found or expired
		}
		return nil, err
	}

	return &refreshToken, nil
}

// RevokeRefreshToken deletes a refresh token from the database
func RevokeRefreshToken(ctx context.Context, token string) error {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := db.ExecContext(dbCtx, query, token)

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Revoke refresh token blocked by circuit breaker", err)
		return errors.New("authentication service temporarily unavailable")
	}

	return err
}

// RevokeAllUserRefreshTokens deletes all refresh tokens for a specific user
func RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := db.ExecContext(dbCtx, query, userID)

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Revoke all user refresh tokens blocked by circuit breaker", err)
		return errors.New("authentication service temporarily unavailable")
	}

	return err
}

// CleanupExpiredTokens removes all expired refresh tokens
func CleanupExpiredTokens(ctx context.Context) error {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM refresh_tokens WHERE expires_at <= NOW() AT TIME ZONE 'UTC'`
	_, err := db.ExecContext(dbCtx, query)

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Cleanup expired tokens blocked by circuit breaker", err)
		return errors.New("authentication service temporarily unavailable")
	}

	return err
}
