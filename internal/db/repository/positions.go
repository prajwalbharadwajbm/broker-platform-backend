package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
)

// GetUserPositions retrieves all positions for a user for PNL calculation
func GetUserPositions(ctx context.Context, userID uuid.UUID) ([]models.Position, error) {
	db := db.GetClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, user_id, symbol, position_type, quantity, entry_price, current_price, unrealized_pnl, realized_pnl, created_at, updated_at 
			  FROM positions 
			  WHERE user_id = $1`

	rows, err := db.QueryContext(dbCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []models.Position

	for rows.Next() {
		var position models.Position
		err := rows.Scan(&position.ID, &position.UserID, &position.Symbol, &position.PositionType,
			&position.Quantity, &position.EntryPrice, &position.CurrentPrice,
			&position.UnrealizedPNL, &position.RealizedPNL, &position.CreatedAt, &position.UpdatedAt)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return positions, nil
}
