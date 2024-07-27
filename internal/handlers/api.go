package handlers

import (
	"net/http"
	// "waas/internal/handlers/keymanagement"
	"waas/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalTransactions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "waas_total_transactions",
			Help: "Total number of transactions processed",
		},
		[]string{"status"},
	)
	transactionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "waas_transaction_duration_seconds",
			Help:    "Duration of transactions",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)
	walletCreation = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "waas_wallet_creation",
			Help: "Total number of wallet creation requests",
		},
		[]string{"status", "walletAddress", "business"}, // Add the new fields here
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(totalTransactions)
	prometheus.MustRegister(transactionDuration)
	prometheus.MustRegister(walletCreation)
}

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

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// r.Route("/getAPIKey", func(router chi.Router) {
	// 	router.Get("/", keymanagement.GenerateAPIKey)
	// })

	r.Route("/create", func(router chi.Router) {
		router.Use(middleware.ValidateAPIKey)
		router.Post("/", CreateWallet)
	})

	r.Route("/send", func(router chi.Router) {
		router.Use(middleware.ValidateAPIKey)
		router.Post("/", SendToken)
	})

	// r.Route("/tokens", func(router chi.Router) {
	// 	router.Get("/", GetSupportedTokens)
	// })

	// r.Route("/wallet/sendtoken", func(router chi.Router) {
	// 	router.Use(middleware.ValidateAPIKey)
	// 	router.Post("/", SendTokens)
	// })
}
