package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/gateway/internal/types"
)

type authHandler struct {
	Client pb.AuthServiceClient
}

func NewAuthHandler(client pb.AuthServiceClient) *authHandler {
	return &authHandler{Client: client}
}

func (h *authHandler) RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/create-user", h.createUser)
	mux.HandleFunc("POST /api/login", h.login)
}

func jsonDecode[T any](r *http.Request, dest *T) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func (h *authHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	if err := jsonDecode(r, &req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	resp, err := h.Client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": resp.Message,
		"success": resp.Success,
	})
}

func (h *authHandler) login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := jsonDecode(r, &req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	resp, err := h.Client.Login(ctx, &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		log.Printf("Error login user: %v", err)
		http.Error(w, "Error login user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": resp.Token,
		"success": resp.Success,
	})
}
