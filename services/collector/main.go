package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	ticket "github.com/waste3d/ai-ops/gen/go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "tickets.new"
)

var writer *kafka.Writer

func main() {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	http.HandleFunc("/ticket", createTicketHandler)
	log.Println("Collector service listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start collector service: %v", err)
	}
}

type CreateTicketRequest struct {
	Source  string `json:"source"`
	Payload string `json:"payload"`
}

func createTicketHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := &ticket.TicketCreatedEvent{
		Id:        uuid.New().String(),
		Source:    req.Source,
		Payload:   req.Payload,
		CreatedAt: timestamppb.New(time.Now()),
	}
	log.Printf("Generated event: %v", event)

	// in binary format
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		http.Error(w, "failed to marshal event", http.StatusInternalServerError)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: eventBytes,
	})
	if err != nil {
		http.Error(w, "failed to write message to kafka", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Ticket accepted: " + event.Id))
}
