package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/waste3d/ai-ops/services/collector/internal/application"
)

type Handler struct {
	useCase *application.TicketUseCase
}

func NewHandler(useCase *application.TicketUseCase) *Handler {
	return &Handler{useCase: useCase}
}

type CreateTicketRequest struct {
	Source  string `json:"source"`
	Payload string `json:"payload"`
}

func (h *Handler) CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := h.useCase.CreateTicket(r.Context(), req.Source, req.Payload)
	if err != nil {
		log.Printf("Failed to create ticket: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Ticket created: %v", ticket)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Ticket accepted: " + ticket.ID))
}
