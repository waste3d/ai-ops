package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	ticket "github.com/waste3d/ai-ops/gen/go"
	"github.com/waste3d/ai-ops/services/auditor/internal/application"
	"github.com/waste3d/ai-ops/services/auditor/internal/domain"
	"google.golang.org/protobuf/proto"
)

type Consumer struct {
	useCase *application.TicketUseCase
	brokers []string
}

func NewConsumer(useCase *application.TicketUseCase, brokers []string) *Consumer {
	return &Consumer{useCase: useCase, brokers: brokers}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumers...")
	var wg sync.WaitGroup

	NewTicketReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          c.brokers,
		GroupID:          "auditor-new-tickets-group-v2",
		Topic:            "tickets.new",
		MaxWait:          500 * time.Millisecond,
		CommitInterval:   time.Second,
		SessionTimeout:   10 * time.Second,
		RebalanceTimeout: 5 * time.Second,
		JoinGroupBackoff: 250 * time.Millisecond,
	})
	AnalyzedTicketReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:          c.brokers,
		GroupID:          "auditor-analyzed-tickets-group-v2",
		Topic:            "tickets.analyzed",
		MaxWait:          500 * time.Millisecond,
		CommitInterval:   time.Second,
		SessionTimeout:   10 * time.Second,
		RebalanceTimeout: 5 * time.Second,
		JoinGroupBackoff: 250 * time.Millisecond,
	})

	wg.Add(2)

	go c.consume(ctx, &wg, NewTicketReader, c.handleNewTicketMessage)
	go c.consume(ctx, &wg, AnalyzedTicketReader, c.handleAnalyzedTicketMessage)

	<-ctx.Done()
	log.Println("Shutdown signal received by consumer, closing readers...")

	NewTicketReader.Close()
	AnalyzedTicketReader.Close()

	wg.Wait()
	log.Println("All consumers have stopped gracefully.")
}

func (c *Consumer) consume(ctx context.Context, wg *sync.WaitGroup, reader *kafka.Reader, handler func(context.Context, kafka.Message)) {
	defer wg.Done()
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Printf("Consumer for topic %s stopped.", reader.Config().Topic)
				return
			}
			log.Printf("Error fetching message from %s: %v", reader.Config().Topic, err)
			continue
		}
		handler(ctx, msg)

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
		}
	}
}

func (c *Consumer) handleNewTicketMessage(ctx context.Context, msg kafka.Message) {
	var event ticket.TicketCreatedEvent
	if err := proto.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Unmarshal error for new ticket: %v", err)
		return
	}

	ticket := &domain.Ticket{
		ID:        event.Id,
		Source:    event.Source,
		Payload:   event.Payload,
		Status:    "new",
		CreatedAt: event.CreatedAt.AsTime(),
	}

	if err := c.useCase.CreateTicket(ctx, ticket); err != nil {
		log.Printf("Failed to create ticket %s: %v", event.Id, err)
		return
	}
	log.Printf("💾 Ticket saved to DB: ID=%s", ticket.ID)
}

func (c *Consumer) handleAnalyzedTicketMessage(ctx context.Context, msg kafka.Message) {
	var event ticket.AnalysisCompletedEvent
	if err := proto.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Unmarshal error for analyzed ticket: %v", err)
		return
	}

	ticket := &domain.Ticket{
		ID:             event.TicketId,
		Status:         "analyzed",
		AnalysisResult: event.Result,
		CreatedAt:      event.AnalyzedAt.AsTime(),
	}

	if err := c.useCase.UpdateTicket(ctx, ticket.ID, ticket.Status, ticket.AnalysisResult); err != nil {
		log.Printf("Failed to update ticket %s: %v", event.TicketId, err)
		return
	}
	log.Printf("✅ Ticket updated in DB: ID=%s", ticket.ID)
}
