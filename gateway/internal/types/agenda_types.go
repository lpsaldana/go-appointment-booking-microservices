package types

type CreateSlotRequest struct {
	ProfessionalID uint   `json:"professional_id"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

type ListAvailableSlotsRequest struct {
	ProfessionalID uint   `json:"professional_id"`
	Date           string `json:"date"`
}

type BookAppointmentRequest struct {
	ClientID uint `json:"client_id"`
	SlotID   uint `json:"slot_id"`
}

type ListAppointmentsRequest struct {
	ClientID       uint `json:"client_id,omitempty"`
	ProfessionalID uint `json:"professional_id,omitempty"`
}
