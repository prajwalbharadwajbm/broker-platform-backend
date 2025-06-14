package repository

import (
	"context"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
)

// GetOrderbookEntries retrieves all orderbook entries
func GetOrderbookEntries(ctx context.Context) ([]models.OrderbookEntry, error) {
	db := db.GetClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, symbol, side, price, quantity, created_at, updated_at 
			  FROM orderbook 
			  ORDER BY symbol, 
			  CASE WHEN side = 'BUY' THEN price END DESC,
			  CASE WHEN side = 'SELL' THEN price END ASC`

	rows, err := db.QueryContext(dbCtx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.OrderbookEntry

	for rows.Next() {
		var entry models.OrderbookEntry
		err := rows.Scan(&entry.ID, &entry.Symbol, &entry.Side, &entry.Price, &entry.Quantity, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetOrderbookEntriesBySymbol retrieves orderbook entries for a specific symbol
func GetOrderbookEntriesBySymbol(ctx context.Context, symbol string) ([]models.OrderbookEntry, error) {
	db := db.GetClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, symbol, side, price, quantity, created_at, updated_at 
			  FROM orderbook 
			  WHERE symbol = $1
			  ORDER BY 
			  CASE WHEN side = 'BUY' THEN price END DESC,
			  CASE WHEN side = 'SELL' THEN price END ASC`

	rows, err := db.QueryContext(dbCtx, query, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.OrderbookEntry

	for rows.Next() {
		var entry models.OrderbookEntry
		err := rows.Scan(&entry.ID, &entry.Symbol, &entry.Side, &entry.Price, &entry.Quantity, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
