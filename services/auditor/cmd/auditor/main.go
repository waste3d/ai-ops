package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	ticketpb "github.com/waste3d/ai-ops/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waste3d/ai-ops/services/auditor/internal/application"
	grpc_infra "github.com/waste3d/ai-ops/services/auditor/internal/infrastructure/grpc"
	"github.com/waste3d/ai-ops/services/auditor/internal/infrastructure/kafka"
	"github.com/waste3d/ai-ops/services/auditor/internal/infrastructure/persistence"
)

const (
	databaseURL = "postgres://user:password@localhost:5432/ops_copilot_db"
	kafkaBroker = "localhost:9092"
	grpcAddr    = ":50051"
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		kafkaConsumer.Start(ctx)
	}()

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := grpc_infra.NewServer(ticketUseCase)
	ticketpb.RegisterAuditServiceServer(grpcServer, grpcHandler)

	reflection.Register(grpcServer)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server failed to serve: %v", err)
		}
	}()

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")

	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()

	cancel()

	wg.Wait()

	dbpool.Close()

	log.Println("Service shutdown complete.")
}
