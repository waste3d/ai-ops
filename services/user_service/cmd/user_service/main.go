package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	authpb "github.com/waste3d/ai-ops/gen/go/auth"
	"github.com/waste3d/ai-ops/services/user_service/internal/application"
	grpc_infra "github.com/waste3d/ai-ops/services/user_service/internal/infrastructure/grpc"
	"github.com/waste3d/ai-ops/services/user_service/internal/infrastructure/persistence"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	databaseURL = "postgres://user:password@localhost:5432/ops_copilot_db"
	grpcAddr    = ":50052"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbpool.Close()
	log.Println("Database connection established for user_service.")

	userRepo := persistence.NewPostgresUserRepository(dbpool)
	userUseCase := application.NewUserUseCase(userRepo)

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcHandler := grpc_infra.NewServer(userUseCase)
	authpb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")
	grpcServer.GracefulStop()

	cancel()

	log.Println("Service shutdown complete.")
}
