package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/prajwalbharadwajbm/broker/internal/db/models"
	"github.com/prajwalbharadwajbm/broker/internal/db/repository"
	"github.com/prajwalbharadwajbm/broker/internal/interceptor"
	"github.com/prajwalbharadwajbm/broker/internal/logger"
	"github.com/prajwalbharadwajbm/broker/internal/utils"
)

func GetHoldings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := ctx.Value("userId").(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.Error("failed to parse user ID", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	holdings, err := repository.GetHoldings(ctx, userUUID)
	if err != nil {
		logger.Log.Error("failed to get holdings", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	interceptor.SendSuccessResponse(w, holdings, http.StatusOK)
}

func AddHolding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := ctx.Value("userId").(string)
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.Error("failed to parse user ID", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	holding, err := utils.FetchDataFromRequestBody[models.Holding](r)
	if err != nil {
		logger.Log.Error("failed to fetch holding from request body", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	holding.UserID = userUUID

	err = repository.AddHolding(ctx, holding)
	if err != nil {
		logger.Log.Error("failed to add holding", err)
		interceptor.SendErrorResponse(w, "BPB009", http.StatusInternalServerError)
		return
	}

	interceptor.SendSuccessResponse(w, "Holding added successfully", http.StatusOK)
}
