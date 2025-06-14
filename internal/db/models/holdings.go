package models

import (
	"time"

	"github.com/google/uuid"
)

// Holding represents a user's holding of a particular asset
// Decimal types are used for financial calculations to ensure precision
type Holding struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	Symbol       string    `json:"symbol" db:"symbol"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	AveragePrice float64   `json:"average_price" db:"average_price"`
	CurrentPrice float64   `json:"current_price" db:"current_price"`
	TotalValue   float64   `json:"total_value" db:"total_value"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
