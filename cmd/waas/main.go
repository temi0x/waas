package main

import (
	"fmt"
	"net/http"
	"waas/internal/database"
	"waas/internal/handlers"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"

	log "github.com/sirupsen/logrus"
)

func main() {
	DB, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to set up the database: %e", err)
	} else {
		log.Printf("Database has started: %v", DB)
	}

	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting Server on 8080")

	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Error("Error starting server", err)
	}

	err = DB.Close()
	if err != nil {
		log.Fatalf("Failed to close DB: %v", err)
		return
	} else {
		log.Info("DB Closed")
	}
}
