package main

import (
	"fmt"
	"net/http"

	"waas/internal/handlers"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)



func main() {
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()
	handlers.Handler(r)

	fmt.Println("Starting Server")

	err := http.ListenAndServe("localhost:8000", r)
	if err != nil {
		log.Error("Error starting server", err)
	}
}
