package main

import (
	"log"
	"net"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/notification/internal/config"
	"github.com/lpsaldana/go-appointment-booking-microservices/notification/internal/handlers"
	"github.com/lpsaldana/go-appointment-booking-microservices/notification/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	smtpConfig := config.NewSMTPConfig()

	clientsConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to client service: %v", err)
	}
	defer clientsConn.Close()

	profConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to profesional service: %v", err)
	}
	defer profConn.Close()

	svc := services.NewNotificationService(smtpConfig, clientsConn, profConn)
	handler := handlers.NewNotificationHandler(svc)

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Error listening to port 50055: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, handler)

	log.Println("Server runing on port :50055...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
