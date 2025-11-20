package grpc

import (
	"context"

	ticketpb "github.com/waste3d/ai-ops/gen/go/ticket"
	"github.com/waste3d/ai-ops/services/auditor/internal/application"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	ticketpb.UnimplementedAuditServiceServer
	useCase *application.TicketUseCase
}

// RegisterService implements grpc.ServiceRegistrar.
func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	panic("unimplemented")
}

func NewServer(useCase *application.TicketUseCase) *Server {
	return &Server{useCase: useCase}
}

func (s *Server) GetAllTickets(ctx context.Context, req *emptypb.Empty) (*ticketpb.GetAllTicketsResponse, error) {
	domainTickets, err := s.useCase.GetAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	pbTickets := make([]*ticketpb.Ticket, 0, len(domainTickets))
	for _, t := range domainTickets {
		pbTickets = append(pbTickets, &ticketpb.Ticket{
			Id:             t.ID,
			Source:         t.Source,
			Payload:        t.Payload,
			Status:         t.Status,
			AnalysisResult: t.AnalysisResult,
			CreatedAt:      timestamppb.New(t.CreatedAt),
		})
	}

	return &ticketpb.GetAllTicketsResponse{Tickets: pbTickets}, nil
}

func (s *Server) GetTicketByID(ctx context.Context, req *ticketpb.GetTicketByIDRequest) (*ticketpb.Ticket, error) {
	domainTicket, err := s.useCase.GetTicketByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	pbTicket := &ticketpb.Ticket{
		Id:             domainTicket.ID,
		Source:         domainTicket.Source,
		Payload:        domainTicket.Payload,
		Status:         domainTicket.Status,
		AnalysisResult: domainTicket.AnalysisResult,
		CreatedAt:      timestamppb.New(domainTicket.CreatedAt),
	}

	return pbTicket, nil
}
