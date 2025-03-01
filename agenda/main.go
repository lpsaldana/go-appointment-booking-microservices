package main

import (
	"log"
	"net"

	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/config"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/handlers"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	dsn = common.EnvString("AGENDA_DB", "host=localhost user=postgres password=yourpassword dbname=agenda_db port=5432 sslmode=disable")
)

func main() {
	dbConfig := config.NewDBConfig(dsn)
	db, err := dbConfig.ConnectDB()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	notifConn, err := grpc.NewClient("localhost:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to notification server: %v", err)
	}
	defer notifConn.Close()

	repo := repositories.NewAgendaRepository(db)
	svc := services.NewAgendaService(repo, notifConn)
	handler := handlers.NewAgendaHandler(svc)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Error listening to port 50054: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAgendaServiceServer(grpcServer, handler)

	log.Println("Server runing on port :50054...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error starting gRPC server: %v", err)
	}
}
