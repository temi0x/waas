package handlers

import (
	"net/http"
	// "waas/internal/handlers/keymanagement"
	"waas/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust this to allow only specific origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Handler(r *chi.Mux) {
	// CORS
	r.Use(enableCORS)
	// Global Middleware
	r.Use(chimiddle.StripSlashes)

	// r.Route("/getAPIKey", func(router chi.Router) {
	// 	router.Get("/", keymanagement.GenerateAPIKey)
	// })

	r.Route("/wallet/create", func(router chi.Router) {
		router.Use(middleware.ValidateAPIKey)
		router.Post("/", CreateWallet)
	})

	r.Route("/wallet/send", func(router chi.Router) {
		router.Use(middleware.ValidateAPIKey)
		router.Post("/", SendToken)
	})

	r.Route("/wallet/sendtoken", func(router chi.Router) {
		router.Use(middleware.ValidateAPIKey)
		router.Post("/", SendCustomToken)
	})
}
