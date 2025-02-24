package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
)

type ClientHandler struct {
	pb.UnimplementedClientServiceServer
	Service services.ClientService
}

func NewClientHandler(svc services.ClientService) *ClientHandler {
	return &ClientHandler{Service: svc}
}

func (h *ClientHandler) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	return h.Service.CreateClient(req)
}

func (h *ClientHandler) GetClient(ctx context.Context, req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	return h.Service.GetClient(req)
}

func (h *ClientHandler) ListClients(ctx context.Context, req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	return h.Service.ListClients(req)
}
