package services

import (
	"context"
	"fmt"
	"log"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/notification/internal/config"
	"google.golang.org/grpc"
)

type NotificationService interface {
	SendAppointmentNotification(clientID, professionalID, appointmentID uint32, startTime, endTime string) (string, bool, error)
}

type NotificationServiceImpl struct {
	SMTPConfig    *config.SMTPConfig
	ClientsClient pb.ClientServiceClient
	ProfClient    pb.ProfessionalServiceClient
}

func NewNotificationService(smtpConfig *config.SMTPConfig, clientsConn, profConn *grpc.ClientConn) NotificationService {
	return &NotificationServiceImpl{
		SMTPConfig:    smtpConfig,
		ClientsClient: pb.NewClientServiceClient(clientsConn),
		ProfClient:    pb.NewProfessionalServiceClient(profConn),
	}
}

func (s *NotificationServiceImpl) SendAppointmentNotification(clientID, professionalID, appointmentID uint32, startTime, endTime string) (string, bool, error) {
	clientResp, err := s.ClientsClient.GetClient(context.TODO(), &pb.GetClientRequest{Id: clientID})
	if err != nil {
		log.Printf("Error obtaining client data: %v", err)
		return "Error obtaining client data", false, err
	}

	profResp, err := s.ProfClient.GetProfessional(context.TODO(), &pb.GetProfessionalRequest{Id: professionalID})
	if err != nil {
		log.Printf("Error obtaining professional data: %v", err)
		return "Error obtaining professional data", false, err
	}

	clientEmail := clientResp.Client.Email
	profEmail := profResp.Professional.Contact

	subject := "Cita Registrada Exitosamente"
	body := fmt.Sprintf("Estimado/a,\n\nSu cita ha sido registrada exitosamente.\n\n"+
		"Detalles de la cita:\n"+
		"- ID de la cita: %d\n"+
		"- Inicio: %s\n"+
		"- Fin: %s\n\n"+
		"Gracias por usar nuestro sistema.\nSaludos,\nEquipo de Agendamiento",
		appointmentID, startTime, endTime)

	err = s.SMTPConfig.SendMail([]string{clientEmail}, subject, body)
	if err != nil {
		return "Error sending client notification", false, err
	}

	err = s.SMTPConfig.SendMail([]string{profEmail}, subject, body)
	if err != nil {
		return "Error sending client notification", false, err
	}

	return "Notification send success", true, nil
}
