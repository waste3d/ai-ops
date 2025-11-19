package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	ticketpb "github.com/waste3d/ai-ops/gen/go"
	"github.com/waste3d/ai-ops/services/collector/internal/application"
	"github.com/waste3d/ai-ops/services/collector/internal/domain"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Producer struct {
	writer *kafka.Writer
}

// checking if Producer implements TicketPublisher interface
var _ application.TicketPublisher = (*Producer)(nil)

func NewKafkaProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTicketCreatedEvent(ctx context.Context, ticket *domain.Ticket) error {
	event := &ticketpb.TicketCreatedEvent{
		Id:        ticket.ID,
		Source:    ticket.Source,
		Payload:   ticket.Payload,
		CreatedAt: timestamppb.New(ticket.CreatedAt),
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(ticket.ID),
		Value: eventBytes,
	})

	if err != nil {
		log.Printf("Failed to write message to kafka: %v", err)
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
