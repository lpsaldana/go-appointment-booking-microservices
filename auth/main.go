package main

import (
	"log"
	"net"

	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/config"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/handlers"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"google.golang.org/grpc"
)

var (
	secretKey = common.EnvString("JWT_SECRET", "please-dont-use-this-key-12345")
	dsn       = common.EnvString("AUTH_DB", "host=localhost user=postgres password=postgres dbname=Auth port=5432 sslmode=disable")
)

func main() {
	dbConfig := config.NewDBConfig(dsn)
	db, err := dbConfig.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	repo := repositories.NewUserRepository(db)
	srv := services.NewAuthService(repo, secretKey)
	handler := handlers.NewAuthHandler(srv)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error opening port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, handler)

	log.Println("Auth server runing in port :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting auth grpc server: %v", err)
	}
}
