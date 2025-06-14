package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/config"
	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/service/auth"
)

const VERSION = "1.0.0"

func init() {
	config.LoadConfigs()
	initializeGlobalLogger()
	loadDatabaseClient()
	logger.Log.Info("loaded all configs")
}

func initializeGlobalLogger() {
	env := config.AppConfigInstance.GeneralConfig.Env
	logLevel := config.AppConfigInstance.GeneralConfig.LogLevel
	logger.InitializeGlobalLogger(logLevel, env, VERSION+"-broker-platform")
	logger.Log.Info("loaded the global logger")
}

func loadDatabaseClient() {
	db.GetClient()
}

func main() {
	// Start token cleanup service in background
	ctx := context.Background()
	// cleanup expired refresh tokens
	go auth.StartTokenCleanupService(ctx)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.AppConfigInstance.GeneralConfig.Port),
		Handler:      Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Log.Infof("Starting server on port %d", config.AppConfigInstance.GeneralConfig.Port)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("failed to serve http server", err)
	}
}
