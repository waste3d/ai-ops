package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
	llm "github.com/waste3d/ai-ops/services/ai_responser/internal/infrastructure"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/infrastructure/kafka"
)

const (
	kafkaBroker     = "localhost:9092"
	inputTopic      = "tickets.new"
	outputTopic     = "tickets.analyzed"
	consumerGroupID = "ai-reasoner-group-v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	llmClient := llm.NewSimulatedClient()
	kafkaProducer := kafka.NewProducer([]string{kafkaBroker}, outputTopic)
	analysisUseCase := application.NewAnalysisUseCase(llmClient, kafkaProducer)
	kafkaConsumer := kafka.NewConsumer(analysisUseCase, []string{kafkaBroker}, inputTopic, consumerGroupID)

	wg.Add(1)
	go kafkaConsumer.Start(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")
	cancel()

	wg.Wait()

	if err := kafkaProducer.Close(); err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
	}

	log.Println("AI Reasoner service stopped.")
}
