package application

import (
	"context"

	authpb "github.com/waste3d/ai-ops/gen/go/auth"
	"github.com/waste3d/ai-ops/services/api_gateway/domain"
)

type AuditorReader interface {
	GetAllTickets(ctx context.Context) ([]*domain.TicketView, error)
	GetTicketByID(ctx context.Context, id string) (*domain.TicketView, error)
}

type TicketUseCase struct {
	auditor AuditorReader
}

type UserCreator interface {
	Register(ctx context.Context, username, passwordHash string) (*authpb.User, error)
}

type UserReader interface {
	GetUserByUsername(ctx context.Context, username string) (*authpb.GetUserByUsernameResponse, error)
}

func NewTicketUseCase(auditor AuditorReader) *TicketUseCase {
	return &TicketUseCase{auditor: auditor}
}

func (uc *TicketUseCase) GetAllTickets(ctx context.Context) ([]*domain.TicketView, error) {
	return uc.auditor.GetAllTickets(ctx)
}

func (uc *TicketUseCase) GetTicketByID(ctx context.Context, id string) (*domain.TicketView, error) {
	return uc.auditor.GetTicketByID(ctx, id)
}
