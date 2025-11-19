package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waste3d/ai-ops/services/auditor/internal/application"
	"github.com/waste3d/ai-ops/services/auditor/internal/infrastructure/kafka"
	"github.com/waste3d/ai-ops/services/auditor/internal/infrastructure/persistence"
)

const (
	databaseURL = "postgres://user:password@localhost:5432/ops_copilot_db"
	kafkaBroker = "localhost:9092"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()
	log.Println("Database connection established.")

	// infra -> persistence -> domain
	ticketRepo := persistence.NewPostgresTicketRepository(dbpool)
	ticketUseCase := application.NewTicketUseCase(ticketRepo)
	kafkaConsumer := kafka.NewConsumer(ticketUseCase, []string{kafkaBroker})

	go kafkaConsumer.Start(ctx)

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := kafkaConsumer.Close(shutdownCtx); err != nil {
		log.Printf("Failed to close Kafka consumers: %v", err)
	}

	dbpool.Close()

	log.Println("Service shutdown complete.")
}
