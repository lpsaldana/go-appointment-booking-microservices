package main

import (
	"log"
	"net"

	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/config"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/handlers"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"google.golang.org/grpc"
)

var (
	dsn = common.EnvString("AUTH_DB", "host=localhost user=postgres password=postgres dbname=clients port=5432 sslmode=disable")
)

func main() {
	dbConfig := config.NewDBConfig(dsn)
	db, err := dbConfig.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	repo := repositories.NewClientRepository(db)
	svc := services.NewClientService(repo)
	handler := handlers.NewClientHandler(svc)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error listening to port 50053: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(grpcServer, handler)

	log.Println("Server runing on port :50053...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
