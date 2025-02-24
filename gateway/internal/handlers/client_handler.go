package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lpsaldana/go-appointment-booking-microservices/common/pb"
	"github.com/lpsaldana/go-appointment-booking-microservices/gateway/internal/middleware"
	"github.com/lpsaldana/go-appointment-booking-microservices/gateway/internal/types"
	"google.golang.org/grpc"
)

type ClientHandler struct {
	Client pb.ClientServiceClient
}

func NewClientHandler(conn *grpc.ClientConn) *ClientHandler {
	return &ClientHandler{Client: pb.NewClientServiceClient(conn)}
}

func (h *ClientHandler) RegisterClientRoutes(mux *http.ServeMux, secretKey string) {
	mux.HandleFunc("POST /api/create-client", middleware.JWTAuthMiddleware(secretKey, h.CreateClientHandler))
	mux.HandleFunc("GET /api/get-client", middleware.JWTAuthMiddleware(secretKey, h.GetClientHandler))
	mux.HandleFunc("GET /api/list-clients", middleware.JWTAuthMiddleware(secretKey, h.ListClientsHandler))
}

func (h *ClientHandler) CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.CreateClient(ctx, &pb.CreateClientRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
	if err != nil {
		http.Error(w, "Error en el servicio de clientes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   resp.Message,
		"success":   resp.Success,
		"client_id": resp.ClientId,
	})
}

func (h *ClientHandler) GetClientHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing required param 'id'", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.GetClient(ctx, &pb.GetClientRequest{
		Id: uint32(id),
	})
	if err != nil {
		http.Error(w, "Error in client service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"client":  resp.Client,
		"success": resp.Success,
	})
}

func (h *ClientHandler) ListClientsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.ListClients(ctx, &pb.ListClientsRequest{})
	if err != nil {
		http.Error(w, "Error in client service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"clients": resp.Clients,
		"success": resp.Success,
	})
}
