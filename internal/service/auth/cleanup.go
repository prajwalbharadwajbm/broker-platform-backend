package auth

import (
	"context"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
)

// StartTokenCleanupService starts a background goroutine to periodically clean up expired refresh tokens
func StartTokenCleanupService(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Run cleanup once per day TODO: make it configurable
	defer ticker.Stop()

	cleanupExpiredTokens(ctx)

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Token cleanup service stopped")
			return
		case <-ticker.C:
			cleanupExpiredTokens(ctx)
		}
	}
}

// cleanupExpiredTokens removes expired refresh tokens from the database
func cleanupExpiredTokens(ctx context.Context) {
	logger.Log.Info("Starting cleanup of expired refresh tokens")

	err := repository.CleanupExpiredTokens(ctx)
	if err != nil {
		logger.Log.Error("Failed to cleanup expired refresh tokens", err)
		return
	}

	logger.Log.Info("Successfully cleaned up expired refresh tokens")
}
