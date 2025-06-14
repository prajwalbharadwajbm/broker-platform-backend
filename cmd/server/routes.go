package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prajwalbharadwajbm/broker/internal/handlers"
	"github.com/prajwalbharadwajbm/broker/internal/middleware"
)

func Routes() *httprouter.Router {
	router := httprouter.New()
	// No auth required for health check endpoint
	router.HandlerFunc(http.MethodGet, "/health", handlers.Health)

	// user endpoints (no auth required)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/signup", handlers.Signup)
	router.HandlerFunc(http.MethodPost, "/api/v1/users/login", handlers.Login)

	// Token refresh endpoints (no auth required)
	router.HandlerFunc(http.MethodPost, "/api/v1/auth/refresh", handlers.RefreshToken)
	router.HandlerFunc(http.MethodPost, "/api/v1/auth/revoke", handlers.RevokeRefreshToken) // Logout

	// authenticated endpoints
	router.HandlerFunc(http.MethodPost, "/api/v1/holdings", middleware.AuthMiddleware(handlers.AddHolding))
	router.HandlerFunc(http.MethodGet, "/api/v1/holdings", middleware.AuthMiddleware(handlers.GetHoldings))

	router.HandlerFunc(http.MethodGet, "/api/v1/orderbook", middleware.AuthMiddleware(handlers.GetOrderbook))
	router.HandlerFunc(http.MethodGet, "/api/v1/positions", middleware.AuthMiddleware(handlers.GetPositions))

	return router
}
