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
	router.HandlerFunc(http.MethodPost, "/api/v1/users/signup", handlers.Signup)

	// authenticated endpoints
	router.HandlerFunc(http.MethodPost, "/api/v1/users/login", handlers.Login)

	return router
}
