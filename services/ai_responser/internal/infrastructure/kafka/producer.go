package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	ticketpb "github.com/waste3d/ai-ops/gen/go/ticket"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/domain"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Producer struct {
	writer *kafka.Writer
}

var _ application.ResultPublisher = (*Producer)(nil)

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishAnalysisResult(ctx context.Context, result *domain.AnalysisResult) error {
	event := &ticketpb.AnalysisCompletedEvent{
		TicketId:   result.TicketID,
		Result:     result.Result,
		AnalyzedAt: timestamppb.New(time.Now()),
	}

	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(result.TicketID),
		Value: eventBytes,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
