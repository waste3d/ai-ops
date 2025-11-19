package main

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
	ticket "github.com/waste3d/ai-ops/gen/go"
	"google.golang.org/protobuf/proto"
)

const (
	kafkaBroker = "localhost:9092"
	databaseURL = "postgres://user:password@localhost:5432/ops_copilot_db"
	topic       = "tickets.new"
	groupID     = "auditor-group"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()

	setupDatabase(dbpool)

	log.Println("Auditor service started, waiting for messages...")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		consumeNewTickets(dbpool)
	}()

	go func() {
		defer wg.Done()
		consumeAnalyzedTickets(dbpool)
	}()

	wg.Wait()
}

func consumeNewTickets(dbpool *pgxpool.Pool) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		GroupID: "auditor-new-tickets-group",
		Topic:   "tickets.new",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("new tickets consumer error: %v", err)
			continue
		}

		var event ticket.TicketCreatedEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("new tickets unmarshal error: %v", err)
			continue
		}

		err = saveTicket(dbpool, &event)
		if err != nil {
			log.Printf("Failed to save ticket %s to DB: %v", event.Id, err)
			continue
		}
		log.Printf("💾 Ticket saved to DB: ID=%s", event.Id)
	}
}

func consumeAnalyzedTickets(dbpool *pgxpool.Pool) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		GroupID: "auditor-analyzed-tickets-group", // Уникальный GroupID
		Topic:   "tickets.analyzed",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("analyzed tickets consumer error: %v", err)
			continue
		}

		var event ticket.AnalysisCompletedEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("analyzed tickets unmarshal error: %v", err)
			continue
		}

		err = updateTicket(dbpool, &event)
		if err != nil {
			log.Printf("Failed to update ticket %s in DB: %v", event.TicketId, err)
			continue
		}
		log.Printf("✅ Ticket updated in DB: ID=%s", event.TicketId)
	}
}

func setupDatabase(dbpool *pgxpool.Pool) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tickets (
		id TEXT PRIMARY KEY,
		source TEXT,
		payload TEXT,
		status TEXT,
		analysis_result TEXT, -- Новая колонка
		created_at TIMESTAMPTZ
	);`
	_, err := dbpool.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v\n", err)
	}
}

func saveTicket(dbpool *pgxpool.Pool, event *ticket.TicketCreatedEvent) error {
	insertSQL := `INSERT INTO tickets (id, source, payload, status, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`
	_, err := dbpool.Exec(context.Background(), insertSQL, event.Id, event.Source, event.Payload, "new", event.GetCreatedAt().AsTime())
	return err
}

// Новая функция для обновления тикета
func updateTicket(dbpool *pgxpool.Pool, event *ticket.AnalysisCompletedEvent) error {
	updateSQL := `UPDATE tickets SET status = $1, analysis_result = $2 WHERE id = $3`
	_, err := dbpool.Exec(context.Background(), updateSQL, "analyzed", event.Result, event.TicketId)
	return err
}
