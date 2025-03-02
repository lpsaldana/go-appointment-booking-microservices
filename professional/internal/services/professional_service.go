package services

import (
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/models"
	"github.com/lpsaldana/go-appointment-booking-microservices/professional/internal/repositories"
)

type ProfessionalService interface {
	CreateProfessional(req *pb.CreateProfessionalRequest) (*pb.CreateProfessionalResponse, error)
	GetProfessional(req *pb.GetProfessionalRequest) (*pb.GetProfessionalResponse, error)
	ListProfessionals(req *pb.ListProfessionalsRequest) (*pb.ListProfessionalsResponse, error)
}

type professionalServiceImpl struct {
	Repo repositories.ProfessionalRepository
}

func NewProfessionalService(repo repositories.ProfessionalRepository) ProfessionalService {
	return &professionalServiceImpl{Repo: repo}
}

func (s *professionalServiceImpl) CreateProfessional(req *pb.CreateProfessionalRequest) (*pb.CreateProfessionalResponse, error) {
	professional := &models.Professional{
		Name:       req.Name,
		Profession: req.Profession,
		Contact:    req.Contact,
	}
	if err := s.Repo.CreateProfessional(professional); err != nil {
		return &pb.CreateProfessionalResponse{
			Message: "Error creating professional",
			Success: false,
		}, err
	}

	return &pb.CreateProfessionalResponse{
		Message:        "Professional created",
		Success:        true,
		ProfessionalId: uint32(professional.ID),
	}, nil
}

func (s *professionalServiceImpl) GetProfessional(req *pb.GetProfessionalRequest) (*pb.GetProfessionalResponse, error) {
	professional, err := s.Repo.GetProfessionalByID(uint(req.Id))
	if err != nil {
		return &pb.GetProfessionalResponse{
			Success: false,
		}, err
	}

	return &pb.GetProfessionalResponse{
		Success: true,
		Professional: &pb.Professional{
			Id:         uint32(professional.ID),
			Name:       professional.Name,
			Profession: professional.Profession,
			Contact:    professional.Contact,
		},
	}, nil
}

func (s *professionalServiceImpl) ListProfessionals(req *pb.ListProfessionalsRequest) (*pb.ListProfessionalsResponse, error) {
	professionals, err := s.Repo.ListProfessionals()
	if err != nil {
		return &pb.ListProfessionalsResponse{
			Success: false,
		}, err
	}

	responseProfessionals := make([]*pb.Professional, len(professionals))

	for i, prof := range professionals {
		responseProfessionals[i] = &pb.Professional{
			Id:         uint32(prof.ID),
			Name:       prof.Name,
			Profession: prof.Profession,
			Contact:    prof.Contact,
		}
	}

	return &pb.ListProfessionalsResponse{
		Professionals: responseProfessionals,
		Success:       true,
	}, nil
}
