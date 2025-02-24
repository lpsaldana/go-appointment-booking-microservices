package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/gateway/internal/handlers"
)

var (
	httpAddr  = common.EnvString("HTTP_ADDR", ":3000")
	secretKey = common.EnvString("JWT_SECRET", "please-dont-use-this-key-12345")
)

func main() {

	mux := http.NewServeMux()

	//auth_server
	authConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to auth service: %v", err)
	}
	defer authConn.Close()

	authClient := pb.NewAuthServiceClient(authConn)
	authHandler := handlers.NewAuthHandler(authClient)
	authHandler.RegisterAuthRoutes(mux)

	//professional_server
	profConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to profesional service: %v", err)
	}
	defer profConn.Close()

	profHandler := handlers.NewProfessionalHandler(profConn)
	profHandler.RegisterProfessionalRoutes(mux, secretKey)

	log.Printf("Starting HTTP server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
