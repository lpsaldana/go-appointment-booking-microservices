package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
)

type AgendaHandler struct {
	pb.UnimplementedAgendaServiceServer
	Service services.AgendaService
}

func NewAgendaHandler(svc services.AgendaService) *AgendaHandler {
	return &AgendaHandler{Service: svc}
}

func (h *AgendaHandler) CreateSlot(ctx context.Context, req *pb.CreateSlotRequest) (*pb.CreateSlotResponse, error) {
	return h.Service.CreateSlot(req)
}

func (h *AgendaHandler) ListAvailableSlots(ctx context.Context, req *pb.ListAvailableSlotsRequest) (*pb.ListAvailableSlotsResponse, error) {
	return h.Service.ListAvailableSlots(req)
}

func (h *AgendaHandler) BookAppointment(ctx context.Context, req *pb.BookAppointmentRequest) (*pb.BookAppointmentResponse, error) {
	return h.Service.BookAppointment(req)
}

func (h *AgendaHandler) ListAppointments(ctx context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	return h.Service.ListAppointments(req)
}
