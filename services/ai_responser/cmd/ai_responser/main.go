package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/waste3d/ai-ops/services/ai_responser/internal/application"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/config"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/infrastructure/kafka"
	"github.com/waste3d/ai-ops/services/ai_responser/internal/infrastructure/llm"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	var llmClient application.Analyzer

	switch cfg.LLM.Provider {
	case "googleai":
		if cfg.LLM.APIKey == "" {
			log.Fatalf("LLM provider is googleai, but API key is not set.")
		}
		log.Println("Using Google AI (Gemini) client for analysis.")

		client, err := llm.NewGoogleAIClient(cfg.LLM.APIKey)
		if err != nil {
			log.Fatalf("Failed to create Google AI client: %v", err)
		}
		llmClient = client
	default:
		log.Fatalf("Unsupported LLM provider: %s", cfg.LLM.Provider)
	}

	kafkaProducer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.Output)
	analysisUseCase := application.NewAnalysisUseCase(llmClient, kafkaProducer)
	kafkaConsumer := kafka.NewConsumer(analysisUseCase, cfg.Kafka.Brokers, cfg.Kafka.Topics.Input, cfg.Kafka.GroupID)

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

	log.Println("AI Responser service stopped.")
}
