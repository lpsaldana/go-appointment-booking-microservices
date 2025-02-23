package handlers

import (
	"context"

	"github.com/lpsaldana/go-appointment-booking-microservices/auth/internal/services"
	"github.com/lpsaldana/go-appointment-booking-microservices/common/api"
)

type AuthHandler struct {
	api.UnimplementedAuthServiceServer
	Service services.AuthService
}

func NewAuthHandler(srv services.AuthService) *AuthHandler {
	return &AuthHandler{Service: srv}
}

func (h *AuthHandler) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return h.Service.CreateUser(req)
}

func (h *AuthHandler) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	return h.Service.Login(req)
}
