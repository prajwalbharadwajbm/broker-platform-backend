package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderbookEntry represents an entry in the orderbook
type OrderbookEntry struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Symbol    string    `json:"symbol" db:"symbol"`
	Side      string    `json:"side" db:"side"`
	Price     float64   `json:"price" db:"price"`
	Quantity  float64   `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
