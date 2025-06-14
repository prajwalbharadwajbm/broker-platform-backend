package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	positions "github.com/prajwalbharadwajbm/broker/internal/service/PNL"
)

// PositionsSummary represents positions summary Card information
type PositionsSummary struct {
	TotalPositions int     `json:"total_positions"`
	LongPositions  int     `json:"long_positions"`
	ShortPositions int     `json:"short_positions"`
	TotalValue     float64 `json:"total_value"`
}

// PositionsResponse represents the complete positions response
type PositionsResponse struct {
	Positions []models.Position `json:"positions"`
	PNL       positions.PNL     `json:"pnl"`
	Summary   PositionsSummary  `json:"summary"`
}

// GetPositions returns all user positions with PNL calculations from the database
func GetPositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("userId").(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.Error("failed to parse user ID", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Processing GET /positions request for user: %s", userId)

	// Fetch user positions from database
	userPositions, err := repository.GetUserPositions(ctx, userUUID)
	if err != nil {
		logger.Log.Error("failed to fetch user positions", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Calculate PNL from user positions
	pnl := positions.NewPNLService().CalculatePositionsPNL(userPositions)

	// Generate summary data
	summary := calculatePositionsSummary(userPositions)

	response := PositionsResponse{
		Positions: userPositions,
		PNL:       pnl,
		Summary:   summary,
	}

	logger.Log.Infof("Successfully fetched %d positions for user: %s", len(userPositions), userId)
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}

// calculatePositionsSummary calculates summary statistics for positions
func calculatePositionsSummary(positions []models.Position) PositionsSummary {
	var summary PositionsSummary
	var totalValue float64

	for _, position := range positions {
		summary.TotalPositions++

		// Count long and short positions
		if position.PositionType == "LONG" {
			summary.LongPositions++
		} else if position.PositionType == "SHORT" {
			summary.ShortPositions++
		}

		// Calculate total value based on current price * quantity
		totalValue += position.CurrentPrice * position.Quantity
	}

	summary.TotalValue = totalValue

	return summary
}
