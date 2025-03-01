package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/notification/internal/services"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	Service services.NotificationService
}

func NewNotificationHandler(svc services.NotificationService) *NotificationHandler {
	return &NotificationHandler{Service: svc}
}

func (h *NotificationHandler) SendAppointmentNotification(ctx context.Context, req *pb.SendAppointmentNotificationRequest) (*pb.SendAppointmentNotificationResponse, error) {
	msg, success, err := h.Service.SendAppointmentNotification(req.ClientId, req.ProfessionalId, req.AppointmentId, req.StartTime, req.EndTime)
	if err != nil {
		return &pb.SendAppointmentNotificationResponse{Message: msg, Success: false}, err
	}
	return &pb.SendAppointmentNotificationResponse{Message: msg, Success: success}, nil
}
