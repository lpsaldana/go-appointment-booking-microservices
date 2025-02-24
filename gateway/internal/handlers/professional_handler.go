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

type ProfessionalHandler struct {
	Client pb.ProfessionalServiceClient
}

func NewProfessionalHandler(conn *grpc.ClientConn) *ProfessionalHandler {
	return &ProfessionalHandler{Client: pb.NewProfessionalServiceClient(conn)}
}

func (h *ProfessionalHandler) RegisterProfessionalRoutes(mux *http.ServeMux, secretKey string) {
	mux.HandleFunc("POST /api/create-professional", middleware.JWTAuthMiddleware(secretKey, h.CreateProfessionalHandler))
	mux.HandleFunc("GET /api/get-professional", middleware.JWTAuthMiddleware(secretKey, h.GetProfessionalHandler))
	mux.HandleFunc("GET /api/list-professionals", middleware.JWTAuthMiddleware(secretKey, h.ListProfessionalsHandler))
}

func (h *ProfessionalHandler) CreateProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateProfessionalRequest
	if err := JsonDecodeInternal(r, &req); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.CreateProfessional(ctx, &pb.CreateProfessionalRequest{
		Name:       req.Name,
		Profession: req.Profession,
		Contact:    req.Contact,
	})
	if err != nil {
		http.Error(w, "Error creating professional", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         resp.Message,
		"success":         resp.Success,
		"professional_id": resp.ProfessionalId,
	})
}

func (h *ProfessionalHandler) GetProfessionalHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id param is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.GetProfessional(ctx, &pb.GetProfessionalRequest{
		Id: uint32(id),
	})
	if err != nil {
		http.Error(w, "Error getting selected professional", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"professional": resp.Professional,
		"success":      resp.Success,
	})
}

func (h *ProfessionalHandler) ListProfessionalsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.ListProfessionals(ctx, &pb.ListProfessionalsRequest{})
	if err != nil {
		http.Error(w, "Error getting professionals list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"professionals": resp.Professionals,
		"success":       resp.Success,
	})
}
