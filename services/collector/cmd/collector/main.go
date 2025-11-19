package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/waste3d/ai-ops/services/collector/internal/application"
	"github.com/waste3d/ai-ops/services/collector/internal/infrastructure"
	inthttp "github.com/waste3d/ai-ops/services/collector/internal/infrastructure/http"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "tickets.new"
	httpAddr    = ":8080"
)

func main() {
	kafkaProducer := infrastructure.NewKafkaProducer([]string{kafkaBroker}, topic)
	ticketUseCase := application.NewTicketUseCase(kafkaProducer)
	httpHandler := inthttp.NewHandler(ticketUseCase)

	server := &http.Server{
		Addr: httpAddr,
	}
	http.HandleFunc("/ticket", httpHandler.CreateTicketHandler)

	go func() {
		log.Printf("Collector service listening on %s", httpAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", httpAddr, err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	<-stopCh

	log.Println("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown gracefully: %v\n", err)
	}

	if err := kafkaProducer.Close(); err != nil {
		log.Printf("Kafka producer shutdown failed: %+v", err)
	}

	log.Println("Collector service stopped.")
}
