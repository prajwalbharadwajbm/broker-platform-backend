package positions

import "github.com/prajwalbharadwajbm/broker/internal/db/models"

// PNL represents PNL summary information
type PNL struct {
	UnrealizedPNL float64 `json:"unrealized_pnl"`
	RealizedPNL   float64 `json:"realized_pnl"`
	TotalPNL      float64 `json:"total_pnl"`
}

// PNLCalculator interface defines methods for calculating PNL
type PNLCalculator interface {
	CalculatePositionsPNL(positions []models.Position) PNL
	CalculateOrderbookPNL(positions []models.Position) PNL

	// TODO: Implement these methods if needed in the future
	CalculatePositionPNL(position models.Position) (unrealized, realized float64)
	CalculateUnrealizedPNL(entryPrice, currentPrice, quantity float64, positionType string) float64
	CalculateRealizedPNL(realizedPNL float64) float64
}

// Service implements PNLCalculator interface
type Service struct{}

// NewPNLService creates a new PNL service instance
func NewPNLService() PNLCalculator {
	return &Service{}
}

// CalculatePositionsPNL calculates PNL summary from user positions
func (s *Service) CalculatePositionsPNL(positions []models.Position) PNL {
	var totalUnrealizedPNL, totalRealizedPNL float64

	for _, position := range positions {
		totalUnrealizedPNL += position.UnrealizedPNL
		totalRealizedPNL += position.RealizedPNL
	}

	return PNL{
		UnrealizedPNL: totalUnrealizedPNL,
		RealizedPNL:   totalRealizedPNL,
		TotalPNL:      totalUnrealizedPNL + totalRealizedPNL,
	}
}

// CalculateOrderbookPNL calculates PNL from positions for orderbook context
// This is essentially the same as CalculatePositionsPNL but kept separate
// if there is any different business logic in the future
func (s *Service) CalculateOrderbookPNL(positions []models.Position) PNL {
	return s.CalculatePositionsPNL(positions)
}
