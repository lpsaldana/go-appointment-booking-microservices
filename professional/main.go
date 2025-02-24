package main

import (
	"log"
	"net"

	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/config"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/handlers"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/services"
	"google.golang.org/grpc"
)

var (
	dsn = common.EnvString("AUTH_DB", "host=localhost user=postgres password=postgres dbname=Professionals port=5432 sslmode=disable")
)

func main() {
	dbConfig := config.NewDBConfig(dsn)
	db, err := dbConfig.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	repo := repositories.NewProfessionalRepository(db)
	svc := services.NewProfessionalService(repo)
	handler := handlers.NewProfessionalHandler(svc)

	lis, err := net.Listen("tcp", ":50052") // Puerto diferente a auth y tasks
	if err != nil {
		log.Fatalf("Error listening to port 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProfessionalServiceServer(grpcServer, handler)

	log.Println("Server runing on port :50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
