package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	Service services.AuthService
}

func NewAuthHandler(srv services.AuthService) *AuthHandler {
	return &AuthHandler{Service: srv}
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return h.Service.CreateUser(req)
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return h.Service.Login(req)
}
