package main

import (
	"context"
	"log"

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

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	log.Println("Auditor service started, waiting for messages...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var event ticket.TicketCreatedEvent
		if err := proto.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("could not unmarshal message: %v", err)
			continue
		}

		log.Printf("✅ Received and processed ticket: ID=%s, Source=%s, Payload='%s'",
			event.Id, event.Source, event.Payload)

		if err := saveTicket(dbpool, &event); err != nil {
			log.Printf("Error saving ticket: %v", err)
			continue
		}

		log.Printf("✅ Ticket saved: ID=%s", event.Id)
	}
}

func setupDatabase(dbpool *pgxpool.Pool) {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tickets (
		id UUID PRIMARY KEY,
		source TEXT,
		payload TEXT,
		status TEXT,
		created_at TIMESTAMPTZ
	);`
	_, err := dbpool.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v\n", err)
	}
}

func saveTicket(dbpool *pgxpool.Pool, event *ticket.TicketCreatedEvent) error {
	createdAt := event.CreatedAt.AsTime()

	insertSQL := `INSERT INTO tickets (id, source, payload, status, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := dbpool.Exec(context.Background(), insertSQL, event.Id, event.Source, event.Payload, "new", createdAt)
	return err
}
