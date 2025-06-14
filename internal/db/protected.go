package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/logger"
	circuit "github.com/rubyist/circuitbreaker"
)

var (
	dbCircuitBreaker *circuit.Breaker
	cbOnce           sync.Once
)

// GetProtectedClient returns a database client protected by circuit breaker
func GetProtectedClient() *ProtectedDB {
	cbOnce.Do(func() {
		// Create a threshold breaker that opens after 5 failures
		// and stays open for 30 seconds before trying again
		dbCircuitBreaker = circuit.NewThresholdBreaker(5)

		logger.Log.Info("Database circuit breaker initialized with threshold: 5")

		// Subscribe to circuit breaker events for monitoring
		go func() {
			events := dbCircuitBreaker.Subscribe()
			for event := range events {
				switch event {
				case circuit.BreakerTripped:
					logger.Log.Error("Database circuit breaker OPENED due to failures", nil)
				case circuit.BreakerReset:
					logger.Log.Info("Database circuit breaker CLOSED - service recovered")
				case circuit.BreakerFail:
					logger.Log.Info("Database operation failed - failure recorded")
				case circuit.BreakerReady:
					logger.Log.Info("Database circuit breaker ready to test recovery")
				}
			}
		}()
	})

	return &ProtectedDB{
		db: GetClient(),
		cb: dbCircuitBreaker,
	}
}

// ProtectedDB wraps the database client with circuit breaker protection
type ProtectedDB struct {
	db *sql.DB
	cb *circuit.Breaker
}

// ExecContext executes a query with circuit breaker protection
func (p *ProtectedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result

	err := p.cb.Call(func() error {
		var err error
		result, err = p.db.ExecContext(ctx, query, args...)
		return err
	}, 5*time.Second) // 5 second timeout

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Database operation blocked - circuit breaker is OPEN", err)
		return nil, err
	}

	return result, err
}

// QueryContext executes a query with circuit breaker protection
func (p *ProtectedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows

	err := p.cb.Call(func() error {
		var err error
		rows, err = p.db.QueryContext(ctx, query, args...)
		return err
	}, 5*time.Second) // 5 second timeout

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Database query blocked - circuit breaker is OPEN", err)
		return nil, err
	}

	return rows, err
}

// QueryRowContext executes a single row query with circuit breaker protection
func (p *ProtectedDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	var row *sql.Row

	err := p.cb.Call(func() error {
		row = p.db.QueryRowContext(ctx, query, args...)
		// Note: QueryRow doesn't return an error until Scan() is called
		// The circuit breaker will handle connection-level errors
		return nil
	}, 5*time.Second) // 5 second timeout

	if err == circuit.ErrBreakerOpen {
		logger.Log.Error("Database query row blocked - circuit breaker is OPEN", err)
		return nil, err
	}

	return row, nil
}

// Ping tests database connectivity with circuit breaker protection
func (p *ProtectedDB) Ping() error {
	return p.cb.Call(func() error {
		return p.db.Ping()
	}, 5*time.Second) // 5 second timeout
}

// GetCircuitBreakerState returns the current state of the circuit breaker for monitoring
func (p *ProtectedDB) GetCircuitBreakerState() string {
	if p.cb.Tripped() {
		return "OPEN"
	}
	return "CLOSED"
}

// GetRawDB returns the underlying database connection (use sparingly)
func (p *ProtectedDB) GetRawDB() *sql.DB {
	return p.db
}

// GetFailures returns the number of failures recorded by the circuit breaker
func (p *ProtectedDB) GetFailures() int64 {
	return p.cb.Failures()
}

// Reset manually resets the circuit breaker (useful for testing or admin operations)
func (p *ProtectedDB) Reset() {
	p.cb.Reset()
	logger.Log.Info("Database circuit breaker manually reset")
}
