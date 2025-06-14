package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	circuit "github.com/rubyist/circuitbreaker"
)

// AddUser demonstrates using circuit breaker for user creation
func AddUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	dbClient := db.GetProtectedClient()

	// Log circuit breaker state for monitoring
	logger.Log.Infof("Database circuit breaker state: %s, failures: %d",
		dbClient.GetCircuitBreakerState(), dbClient.GetFailures())

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	var userId string

	row, err := dbClient.QueryRowContext(dbCtx, query, email, hashedPassword)
	if err != nil {
		// Handle circuit breaker specific errors
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("User creation blocked by circuit breaker", err)
			return "", errors.New("database service temporarily unavailable")
		}
		return "", err
	}

	err = row.Scan(&userId)
	if err != nil {
		return "", err
	}

	return userId, nil
}

// GetUserByEmail demonstrates using circuit breaker for user lookup
func GetUserByEmail(ctx context.Context, email string) (string, []byte, error) {
	dbClient := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userId string
	var hashedPassword []byte

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	row, err := dbClient.QueryRowContext(dbCtx, query, email)
	if err != nil {
		// Handle circuit breaker specific errors
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("User lookup blocked by circuit breaker", err)
			return "", nil, errors.New("database service temporarily unavailable")
		}
		return "", nil, err
	}

	err = row.Scan(&userId, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}

	return userId, hashedPassword, nil
}
