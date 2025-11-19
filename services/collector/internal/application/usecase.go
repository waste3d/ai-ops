package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/waste3d/ai-ops/services/collector/internal/domain"
)

type TicketUseCase struct {
	publisher TicketPublisher
}

func NewTicketUseCase(publisher TicketPublisher) *TicketUseCase {
	return &TicketUseCase{publisher: publisher}
}

func (uc *TicketUseCase) CreateTicket(ctx context.Context, source, payload string) (*domain.Ticket, error) {
	ticket := &domain.Ticket{
		ID:        uuid.New().String(),
		Source:    source,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	if err := uc.publisher.PublishTicketCreatedEvent(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}
