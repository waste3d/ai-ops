package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/waste3d/ai-ops/services/api_gateway/application"
	grpc_client "github.com/waste3d/ai-ops/services/api_gateway/infrastructure/gprc"
	http_handler "github.com/waste3d/ai-ops/services/api_gateway/infrastructure/http"
)

const (
	auditorServiceAddr = "localhost:50051"
	httpAddr           = ":8000"
)

func main() {
	auditorClient, err := grpc_client.NewAuditorClient(context.Background(), auditorServiceAddr)
	if err != nil {
		log.Fatalf("Failed to connect to auditor service: %v", err)
	}
	defer auditorClient.Close()

	ticketUseCase := application.NewTicketUseCase(auditorClient)
	httpHandler := http_handler.NewHandler(ticketUseCase)

	router := gin.Default()
	httpHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	go func() {
		log.Printf("API Gateway listening on %s", httpAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", httpAddr, err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Println("Shutdown signal received, shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %+v", err)
	}

	log.Println("API Gateway stopped.")
}
