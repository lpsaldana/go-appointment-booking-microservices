package services

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/client/internal/repositories"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
)

type ClientService interface {
	CreateClient(req *pb.CreateClientRequest) (*pb.CreateClientResponse, error)
	GetClient(req *pb.GetClientRequest) (*pb.GetClientResponse, error)
	ListClients(req *pb.ListClientsRequest) (*pb.ListClientsResponse, error)
}

type ClientServiceImpl struct {
	Repo repositories.ClientRepository
}

func NewClientService(repo repositories.ClientRepository) ClientService {
	return &ClientServiceImpl{Repo: repo}
}

func (s *ClientServiceImpl) CreateClient(req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	client := &models.Client{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}
	if err := s.Repo.CreateClient(client); err != nil {
		return &pb.CreateClientResponse{
			Message: "Error creating cliente",
			Success: false,
		}, err
	}

	return &pb.CreateClientResponse{
		Message:  "Client created",
		Success:  true,
		ClientId: uint32(client.ID),
	}, nil
}

func (s *ClientServiceImpl) GetClient(req *pb.GetClientRequest) (*pb.GetClientResponse, error) {
	client, err := s.Repo.GetClientByID(uint(req.Id))
	if err != nil {
		return &pb.GetClientResponse{
			Success: false,
		}, err
	}

	return &pb.GetClientResponse{
		Client: &pb.Client{
			Id:    uint32(client.ID),
			Name:  client.Name,
			Email: client.Email,
			Phone: client.Phone,
		},
		Success: true,
	}, nil
}

func (s *ClientServiceImpl) ListClients(req *pb.ListClientsRequest) (*pb.ListClientsResponse, error) {
	clients, err := s.Repo.ListClients()
	if err != nil {
		return &pb.ListClientsResponse{
			Success: false,
		}, err
	}

	pbClients := make([]*pb.Client, len(clients))
	for i, client := range clients {
		pbClients[i] = &pb.Client{
			Id:    uint32(client.ID),
			Name:  client.Name,
			Email: client.Email,
			Phone: client.Phone,
		}
	}

	return &pb.ListClientsResponse{
		Clients: pbClients,
		Success: true,
	}, nil
}
