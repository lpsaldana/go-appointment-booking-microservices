package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/gateway/internal/handlers"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":3000")
)

func main() {
	mux := http.NewServeMux()
	handler := handlers.NewHandler()
	handler.RegisterRoutes(mux)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
