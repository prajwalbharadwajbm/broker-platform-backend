package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	orderbook "github.com/prajwalbharadwajbm/broker/internal/service/PNL"
)

// OrderbookSummary represents orderbook summary statistics
type OrderbookSummary struct {
	TotalBuyOrders  int     `json:"total_buy_orders"`
	TotalSellOrders int     `json:"total_sell_orders"`
	TotalVolume     float64 `json:"total_volume"`
	BestBidPrice    float64 `json:"best_bid_price"`
	BestAskPrice    float64 `json:"best_ask_price"`
	Spread          float64 `json:"spread"`
}

// OrderbookResponse represents the complete orderbook response
type OrderbookResponse struct {
	Entries []models.OrderbookEntry `json:"entries"`
	PNL     orderbook.PNL           `json:"pnl"`
	Summary OrderbookSummary        `json:"summary"`
}

// GetOrderbook returns orderbook data with PNL calculations from the database
func GetOrderbook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID from context (set by auth middleware)
	userId := ctx.Value("userId").(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.Error("failed to parse user ID", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Processing GET /orderbook request for user: %s", userId)

	// Fetch orderbook entries from database
	orderbookEntries, err := repository.GetOrderbookEntries(ctx)
	if err != nil {
		logger.Log.Error("failed to fetch orderbook entries", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Fetch user positions for PNL calculation
	userPositions, err := repository.GetUserPositions(ctx, userUUID)
	if err != nil {
		logger.Log.Error("failed to fetch user positions", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	// Calculate PNL from user positions
	pnl := orderbook.NewPNLService().CalculateOrderbookPNL(userPositions)

	// Generate summary data
	summary := calculateOrderbookSummary(orderbookEntries)

	response := OrderbookResponse{
		Entries: orderbookEntries,
		PNL:     pnl,
		Summary: summary,
	}

	logger.Log.Infof("Successfully fetched orderbook data with %d entries for user: %s", len(orderbookEntries), userId)
	interceptor.SendSuccessResponse(w, response, http.StatusOK)
}

// calculateOrderbookSummary calculates summary statistics for the orderbook
func calculateOrderbookSummary(entries []models.OrderbookEntry) OrderbookSummary {
	var summary OrderbookSummary
	var bestBid, bestAsk float64
	bestAsk = 999999999.0 // Initialize to high value to find minimum

	for _, entry := range entries {
		summary.TotalVolume += entry.Quantity

		if entry.Side == "BUY" {
			summary.TotalBuyOrders++
			if entry.Price > bestBid {
				bestBid = entry.Price
			}
		} else if entry.Side == "SELL" {
			summary.TotalSellOrders++
			if entry.Price < bestAsk {
				bestAsk = entry.Price
			}
		}
	}

	summary.BestBidPrice = bestBid
	summary.BestAskPrice = bestAsk

	// Calculate spread only if we have both bid and ask
	if bestBid > 0 && bestAsk < 999999999.0 {
		summary.Spread = bestAsk - bestBid
	}

	return summary
}
