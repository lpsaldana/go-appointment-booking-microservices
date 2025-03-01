package services

import (
	"time"

	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/agenda/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
)

type AgendaService interface {
	CreateSlot(req *pb.CreateSlotRequest) (*pb.CreateSlotResponse, error)
	ListAvailableSlots(req *pb.ListAvailableSlotsRequest) (*pb.ListAvailableSlotsResponse, error)
	BookAppointment(req *pb.BookAppointmentRequest) (*pb.BookAppointmentResponse, error)
	ListAppointments(req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error)
}

type agendaServiceImpl struct {
	Repo repositories.AgendaRepository
}

func NewAgendaService(repo repositories.AgendaRepository) AgendaService {
	return &agendaServiceImpl{Repo: repo}
}

func (s *agendaServiceImpl) CreateSlot(req *pb.CreateSlotRequest) (*pb.CreateSlotResponse, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return &pb.CreateSlotResponse{Message: "start_time invalid format", Success: false}, err
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return &pb.CreateSlotResponse{Message: "end_time invalid format", Success: false}, err
	}

	slot := &models.Slot{
		ProfessionalID: uint(req.ProfessionalId),
		StartTime:      startTime,
		EndTime:        endTime,
		Available:      true,
	}
	if err := s.Repo.CreateSlot(slot); err != nil {
		return &pb.CreateSlotResponse{Message: "Error creating slot", Success: false}, err
	}

	return &pb.CreateSlotResponse{
		Message: "Slot created",
		Success: true,
		SlotId:  uint32(slot.ID),
	}, nil
}

func (s *agendaServiceImpl) ListAvailableSlots(req *pb.ListAvailableSlotsRequest) (*pb.ListAvailableSlotsResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &pb.ListAvailableSlotsResponse{Success: false}, err
	}

	slots, err := s.Repo.ListAvailableSlots(uint(req.ProfessionalId), date)
	if err != nil {
		return &pb.ListAvailableSlotsResponse{Success: false}, err
	}

	pbSlots := make([]*pb.Slot, len(slots))
	for i, slot := range slots {
		pbSlots[i] = &pb.Slot{
			Id:             uint32(slot.ID),
			ProfessionalId: uint32(slot.ProfessionalID),
			StartTime:      slot.StartTime.Format(time.RFC3339),
			EndTime:        slot.EndTime.Format(time.RFC3339),
			Available:      slot.Available,
		}
	}

	return &pb.ListAvailableSlotsResponse{
		Slots:   pbSlots,
		Success: true,
	}, nil
}

func (s *agendaServiceImpl) BookAppointment(req *pb.BookAppointmentRequest) (*pb.BookAppointmentResponse, error) {
	// Verificar si el slot está disponible usando GetSlotByID
	slot, err := s.Repo.GetSlotByID(uint(req.SlotId))
	if err != nil {
		return &pb.BookAppointmentResponse{Message: "Slot not found", Success: false}, err
	}
	if !slot.Available {
		return &pb.BookAppointmentResponse{Message: "This slot is not available", Success: false}, nil
	}

	appointment := &models.Appointment{
		ClientID:       uint(req.ClientId),
		SlotID:         uint(req.SlotId),
		ProfessionalID: slot.ProfessionalID,
	}
	if err := s.Repo.CreateAppointment(appointment); err != nil {
		return &pb.BookAppointmentResponse{Message: "Error generating appointment", Success: false}, err
	}

	// Marcar el slot como no disponible
	if err := s.Repo.UpdateSlotAvailability(uint(req.SlotId), false); err != nil {
		return &pb.BookAppointmentResponse{Message: "Error updating slot", Success: false}, err
	}

	return &pb.BookAppointmentResponse{
		Message:       "Appointment successfully generated",
		Success:       true,
		AppointmentId: uint32(appointment.ID),
	}, nil
}

func (s *agendaServiceImpl) ListAppointments(req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appointments, err := s.Repo.ListAppointments(uint(req.ClientId), uint(req.ProfessionalId))
	if err != nil {
		return &pb.ListAppointmentsResponse{Success: false}, err
	}

	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, appt := range appointments {
		// Obtener el slot asociado usando GetSlotByID
		slot, err := s.Repo.GetSlotByID(appt.SlotID)
		if err != nil {
			return &pb.ListAppointmentsResponse{Success: false}, err
		}
		pbAppointments[i] = &pb.Appointment{
			Id:             uint32(appt.ID),
			ClientId:       uint32(appt.ClientID),
			SlotId:         uint32(appt.SlotID),
			StartTime:      slot.StartTime.Format(time.RFC3339),
			EndTime:        slot.EndTime.Format(time.RFC3339),
			ProfessionalId: uint32(appt.ProfessionalID),
		}
	}

	return &pb.ListAppointmentsResponse{
		Appointments: pbAppointments,
		Success:      true,
	}, nil
}
