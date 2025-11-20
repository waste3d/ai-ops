package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/auditor/internal/domain"
)

type TicketRepository interface {
	Save(ctx context.Context, ticket *domain.Ticket) error
	Update(ctx context.Context, ticketID string, status string, result string) error
	GetAll(ctx context.Context) ([]*domain.Ticket, error)
	GetByID(ctx context.Context, id string) (*domain.Ticket, error)
}
