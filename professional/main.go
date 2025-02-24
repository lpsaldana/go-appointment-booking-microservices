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
		log.Fatalf("No se pudo inicializar la base de datos: %v", err)
	}

	repo := repositories.NewProfessionalRepository(db)
	svc := services.NewProfessionalService(repo)
	handler := handlers.NewProfessionalHandler(svc)

	lis, err := net.Listen("tcp", ":50052") // Puerto diferente a auth y tasks
	if err != nil {
		log.Fatalf("Error al escuchar en el puerto 50052: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProfessionalServiceServer(grpcServer, handler)

	log.Println("Servidor de profesionales corriendo en :50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al iniciar el servidor gRPC: %v", err)
	}
}
