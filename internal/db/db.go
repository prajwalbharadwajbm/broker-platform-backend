package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/prajwalbharadwajbm/broker/internal/config"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
)

const pgConnStrFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

var (
	client *sql.DB
	once   sync.Once
)

// GetClient returns the database client, initializing it if necessary
func GetClient() *sql.DB {
	once.Do(func() {
		initializeClient()
	})
	return client
}

// Initialize the database client
func initializeClient() {
	connStr := buildConnString()

	var err error
	client, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Log.Fatal("failed to connect to db", err)
	}

	configureConnPoolParams(client)

	// Should Ping to understand if database connection was successful
	err = client.Ping()
	if err != nil {
		logger.Log.Fatal("failed to ping db", err)
	}
	logger.Log.Info("connected to database")
}

func buildConnString() string {
	return fmt.Sprintf(
		pgConnStrFormat,
		config.AppConfigInstance.DB.Host,
		config.AppConfigInstance.DB.Port,
		config.AppConfigInstance.DB.User,
		config.AppConfigInstance.DB.Password,
		config.AppConfigInstance.DB.DBname,
	)
}

// Enabling Connection Pooing and also keeping it to values which runs in local env
// TODO: make them env variable later for prod.

// Benefits:
// 1. Reuses existing connections hence incresing performace of the system, as we are not creating new conn's.
// 2. Prevents database server overload by limiting concurrent connections
// 3. Prevents connection exhaustion errors during high traffic periods
func configureConnPoolParams(client *sql.DB) {
	client.SetMaxOpenConns(25)
	client.SetMaxIdleConns(10)
	client.SetConnMaxLifetime(5 * time.Minute)
	client.SetConnMaxIdleTime(3 * time.Minute)
}
