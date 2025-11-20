package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/auditor/internal/domain"
)

type TicketUseCase struct {
	repo TicketRepository
}

func NewTicketUseCase(repo TicketRepository) *TicketUseCase {
	return &TicketUseCase{repo: repo}
}

func (uc *TicketUseCase) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	// дальше - валидация, логиривация, нотификация
	return uc.repo.Save(ctx, ticket)
}

func (uc *TicketUseCase) UpdateTicket(ctx context.Context, ticketID string, status string, result string) error {
	// дальше - валидация, логиривация, нотификация
	return uc.repo.Update(ctx, ticketID, status, result)
}

func (uc *TicketUseCase) GetAllTickets(ctx context.Context) ([]*domain.Ticket, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *TicketUseCase) GetTicketByID(ctx context.Context, id string) (*domain.Ticket, error) {
	return uc.repo.GetByID(ctx, id)
}
