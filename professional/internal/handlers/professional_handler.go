package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/services"
)

type ProfessionalHandler struct {
	pb.UnimplementedProfessionalServiceServer
	Service services.ProfessionalService
}

func NewProfessionalHandler(service services.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{Service: service}
}

func (h *ProfessionalHandler) CreateProfessional(ctx context.Context, req *pb.CreateProfessionalRequest) (*pb.CreateProfessionalResponse, error) {
	return h.Service.CreateProfessional(req)
}

func (h *ProfessionalHandler) ListProfessionals(ctx context.Context, req *pb.ListProfessionalsRequest) (*pb.ListProfessionalsResponse, error) {
	return h.Service.ListProfessionals(req)
}

func (h *ProfessionalHandler) GetProfessional(ctx context.Context, req *pb.GetProfessionalRequest) (*pb.GetProfessionalResponse, error) {
	return h.Service.GetProfessional(req)
}
