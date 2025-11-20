package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	ticketpb "github.com/waste3d/ai-ops/gen/go"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/domain"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	useCase *application.AnalysisUseCase
	reader  *kafka.Reader
}

func NewConsumer(useCase *application.AnalysisUseCase, brokers []string, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          brokers,
		GroupID:          groupID,
		Topic:            topic,
		MaxWait:          500 * time.Millisecond,
		CommitInterval:   time.Second,
		SessionTimeout:   10 * time.Second,
		RebalanceTimeout: 5 * time.Second,
		JoinGroupBackoff: 250 * time.Millisecond,
	})
	return &Consumer{useCase: useCase, reader: reader}
}

func (c *Consumer) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Printf("Consumer for topic %s stopped.", c.reader.Config().Topic)
				return
			}
			log.Printf("Error fetching message from %s: %v", c.reader.Config().Topic, err)
			continue
		}

		var event ticketpb.TicketCreatedEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Unmarshal error for analyzed ticket: %v", err)
			continue
		}
		log.Printf("🧠 Received ticket for analysis: ID=%s", event.Id)

		ticket := &domain.Ticket{
			ID:      event.Id,
			Payload: event.Payload,
		}

		if err := c.useCase.AnalyzeTicket(ctx, ticket); err != nil {
			log.Printf("Failed to analyze ticket %s: %v", event.Id, err)
			continue
		}
		log.Printf("✅ Published analysis for ticket: ID=%s", ticket.ID)

		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
