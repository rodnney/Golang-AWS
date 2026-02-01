package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/handlers"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/internal/services"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

func main() {
	// 1. Load config and logger
	log := logger.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. Initialize AWS Clients
	awsCfg, err := aws.NewAWSSession(cfg)
	if err != nil {
		log.Error("Failed to initialize AWS session", "error", err)
		os.Exit(1)
	}

	snsClient := aws.NewSNSClient(awsCfg)
	dynamoClient := aws.NewDynamoDBClient(awsCfg)

	// 3. Initialize Repositories and Services
	txRepo := repository.NewDynamoDBTransactionRepository(dynamoClient, cfg)
	txService := services.NewTransactionService(txRepo, snsClient, cfg, log)

	// 4. Initialize Handlers and Router
	txHandler := handlers.NewTransactionHandler(txService, log)
	router := handlers.SetupRouter(txHandler, log)

	// 5. Start HTTP Server
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Starting API server", "addr", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited gracefully")
}
