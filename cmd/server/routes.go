package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prajwalbharadwajbm/broker/internal/handlers"
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

	return router
}
