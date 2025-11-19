package application

import (
	"context"

	"github.com/waste3d/ai-ops/services/collector/internal/domain"
)

type TicketPublisher interface {
	PublishTicketCreatedEvent(ctx context.Context, ticket *domain.Ticket) error
}
