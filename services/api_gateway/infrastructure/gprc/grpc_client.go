package grpc_client

import (
	"context"

	ticketpb "github.com/waste3d/ai-ops/gen/go/ticket"
	"github.com/waste3d/ai-ops/services/api_gateway/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuditorClient struct {
	client ticketpb.AuditServiceClient
	conn   *grpc.ClientConn
}

func NewAuditorClient(ctx context.Context, addr string) (*AuditorClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := ticketpb.NewAuditServiceClient(conn)
	return &AuditorClient{client: client, conn: conn}, nil
}

func (c *AuditorClient) Close() {
	c.conn.Close()
}

func (c *AuditorClient) GetAllTickets(ctx context.Context) ([]*domain.TicketView, error) {
	resp, err := c.client.GetAllTickets(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	tickets := make([]*domain.TicketView, 0, len(resp.Tickets))
	for _, t := range resp.Tickets {
		tickets = append(tickets, &domain.TicketView{
			ID:             t.Id,
			Source:         t.Source,
			Payload:        t.Payload,
			Status:         t.Status,
			AnalysisResult: t.AnalysisResult,
			CreatedAt:      t.CreatedAt.AsTime(),
		})
	}
	return tickets, nil
}

func (c *AuditorClient) GetTicketByID(ctx context.Context, id string) (*domain.TicketView, error) {
	resp, err := c.client.GetTicketByID(ctx, &ticketpb.GetTicketByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}

	ticket := &domain.TicketView{
		ID:             resp.Id,
		Source:         resp.Source,
		Payload:        resp.Payload,
		Status:         resp.Status,
		AnalysisResult: resp.AnalysisResult,
		CreatedAt:      resp.CreatedAt.AsTime(),
	}
	return ticket, nil
}
