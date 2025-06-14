package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	circuit "github.com/rubyist/circuitbreaker"
)

func GetHoldings(ctx context.Context, userId uuid.UUID) ([]models.Holding, error) {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT * FROM holdings WHERE user_id=$1`
	rows, err := db.QueryContext(dbCtx, query, userId)
	if err != nil {
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("Holdings lookup blocked by circuit breaker", err)
			return nil, errors.New("database service temporarily unavailable")
		}
		return nil, err
	}
	defer rows.Close()

	holdings := []models.Holding{}

	for rows.Next() {
		var holding models.Holding
		err := rows.Scan(&holding.ID, &holding.UserID, &holding.Symbol, &holding.Quantity, &holding.AveragePrice, &holding.CurrentPrice, &holding.TotalValue, &holding.CreatedAt, &holding.UpdatedAt)
		if err != nil {
			return nil, err
		}
		holdings = append(holdings, holding)
	}

	return holdings, nil
}

func AddHolding(ctx context.Context, holding models.Holding) error {
	db := db.GetProtectedClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO holdings (user_id, symbol, quantity, average_price, current_price, total_value) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(dbCtx, query, holding.UserID, holding.Symbol, holding.Quantity, holding.AveragePrice, holding.CurrentPrice, holding.TotalValue)
	if err != nil {
		if err == circuit.ErrBreakerOpen {
			logger.Log.Error("Holdings lookup blocked by circuit breaker", err)
			return errors.New("database service temporarily unavailable")
		}
		return err
	}

	return nil
}
