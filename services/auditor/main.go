package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	ticket "github.com/waste3d/ai-ops/gen/go"
	"google.golang.org/protobuf/proto"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "tickets.new"
	groupID     = "auditor-group"
)

func main() {
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
	}
}
