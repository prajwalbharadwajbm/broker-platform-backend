package models

import (
	"time"

	"github.com/google/uuid"
)

// Position represents a user's trading position
// Decimal types are used for financial calculations to ensure precision
type Position struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Symbol        string    `json:"symbol" db:"symbol"`
	PositionType  string    `json:"position_type" db:"position_type"`
	Quantity      float64   `json:"quantity" db:"quantity"`
	EntryPrice    float64   `json:"entry_price" db:"entry_price"`
	CurrentPrice  float64   `json:"current_price" db:"current_price"`
	UnrealizedPNL float64   `json:"unrealized_pnl" db:"unrealized_pnl"`
	RealizedPNL   float64   `json:"realized_pnl" db:"realized_pnl"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
