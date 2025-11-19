package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	ticket "github.com/waste3d/ai-ops/gen/go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	kafkaBroker     = "localhost:9092"
	inputTopic      = "tickets.new"
	outputTopic     = "tickets.analyzed"
	consumerGroupID = "ai-reasoner-group"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   inputTopic,
		GroupID: consumerGroupID,
	})
	defer reader.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    outputTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	log.Println("AI reasoner service started, waiting for messages...")

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

		log.Printf("🧠 Received ticket for analysis: ID=%s", event.Id)

		// Simulate analysis process
		var analysisResult string

		if strings.Contains(strings.ToLower(event.Payload), "баз") {
			analysisResult = "Проблема классифицирована как связанная с базой данных."
		} else if strings.Contains(strings.ToLower(event.Payload), "диск") {
			analysisResult = "Проблема классифицирована как связанная с дисковым пространством."
		} else {
			analysisResult = "Общая проблема, требуется ручной разбор."
		}

		log.Printf("🧠 Analysis result: %s", analysisResult)

		analysisEvent := &ticket.AnalysisCompletedEvent{
			TicketId:   event.Id,
			Result:     analysisResult,
			AnalyzedAt: timestamppb.New(time.Now()),
		}

		eventBytes, err := proto.Marshal(analysisEvent)
		if err != nil {
			log.Printf("could not marshal analysis event: %v", err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Value: eventBytes,
		})
		if err != nil {
			log.Printf("could not write analysis event: %v", err)
			continue
		}

		log.Printf("✅ Published analysis for ticket: ID=%s", analysisEvent.TicketId)

	}
}
