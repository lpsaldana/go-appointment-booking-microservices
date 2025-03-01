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

type AgendaHandler struct {
	Client pb.AgendaServiceClient
}

func NewAgendaHandler(conn *grpc.ClientConn) *AgendaHandler {
	return &AgendaHandler{Client: pb.NewAgendaServiceClient(conn)}
}

func (h *AgendaHandler) RegisterAgendaRoutes(mux *http.ServeMux, secretKey string) {
	mux.HandleFunc("POST /api/create-slot", middleware.JWTAuthMiddleware(secretKey, h.CreateSlotHandler))
	mux.HandleFunc("GET /api/list-available-slots", middleware.JWTAuthMiddleware(secretKey, h.ListAvailableSlotsHandler))
	mux.HandleFunc("POST /api/book-appointment", middleware.JWTAuthMiddleware(secretKey, h.BookAppointmentHandler))
	mux.HandleFunc("GET /api/list-appointments", middleware.JWTAuthMiddleware(secretKey, h.ListAppointmentsHandler))

}

func (h *AgendaHandler) CreateSlotHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.CreateSlot(ctx, &pb.CreateSlotRequest{
		ProfessionalId: uint32(req.ProfessionalID),
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})
	if err != nil {
		http.Error(w, "Error en el servicio de agenda", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": resp.Message,
		"success": resp.Success,
		"slot_id": resp.SlotId,
	})
}

func (h *AgendaHandler) ListAvailableSlotsHandler(w http.ResponseWriter, r *http.Request) {
	profIDStr := r.URL.Query().Get("professional_id")
	date := r.URL.Query().Get("date")
	if profIDStr == "" || date == "" {
		http.Error(w, "Faltan parámetros 'professional_id' o 'date'", http.StatusBadRequest)
		return
	}

	profID, err := strconv.ParseUint(profIDStr, 10, 32)
	if err != nil {
		http.Error(w, "professional_id inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.ListAvailableSlots(ctx, &pb.ListAvailableSlotsRequest{
		ProfessionalId: uint32(profID),
		Date:           date,
	})
	if err != nil {
		http.Error(w, "Error en el servicio de agenda", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"slots":   resp.Slots,
		"success": resp.Success,
	})
}

func (h *AgendaHandler) BookAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	var req types.BookAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.BookAppointment(ctx, &pb.BookAppointmentRequest{
		ClientId: uint32(req.ClientID),
		SlotId:   uint32(req.SlotID),
	})
	if err != nil {
		http.Error(w, "Error en el servicio de agenda", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        resp.Message,
		"success":        resp.Success,
		"appointment_id": resp.AppointmentId,
	})
}

func (h *AgendaHandler) ListAppointmentsHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("client_id")
	profIDStr := r.URL.Query().Get("professional_id")

	var clientID, profID uint32
	if clientIDStr != "" {
		id, err := strconv.ParseUint(clientIDStr, 10, 32)
		if err != nil {
			http.Error(w, "client_id inválido", http.StatusBadRequest)
			return
		}
		clientID = uint32(id)
	}
	if profIDStr != "" {
		id, err := strconv.ParseUint(profIDStr, 10, 32)
		if err != nil {
			http.Error(w, "professional_id inválido", http.StatusBadRequest)
			return
		}
		profID = uint32(id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := h.Client.ListAppointments(ctx, &pb.ListAppointmentsRequest{
		ClientId:       clientID,
		ProfessionalId: profID,
	})
	if err != nil {
		http.Error(w, "Error en el servicio de agenda", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"appointments": resp.Appointments,
		"success":      resp.Success,
	})
}
