package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/internal/services"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	awsCfg, err := aws.NewAWSSession(cfg)
	if err != nil {
		log.Error("Failed to initialize AWS session", "error", err)
		os.Exit(1)
	}

	snsClient := aws.NewSNSClient(awsCfg)
	sqsClient := aws.NewSQSClient(awsCfg)
	dynamoClient := aws.NewDynamoDBClient(awsCfg)
	s3Client := aws.NewS3Client(awsCfg)

	txRepo := repository.NewDynamoDBTransactionRepository(dynamoClient, cfg)
	auditRepo := repository.NewS3AuditRepository(s3Client, cfg)

	validatorSvc := services.NewValidatorService(txRepo, snsClient, cfg, log)
	enrichmentSvc := services.NewEnrichmentService(txRepo, log)
	auditSvc := services.NewAuditService(auditRepo, log)

	var wg sync.WaitGroup
	quit := make(chan struct{})

	// Initialize Workers
	validatorWorker := NewWorker(sqsClient, cfg.SQS.ValidationURL, validatorSvc.Process, log, &wg, quit)
	enrichmentWorker := NewWorker(sqsClient, cfg.SQS.EnrichmentURL, enrichmentSvc.Process, log, &wg, quit)
	auditWorker := NewWorker(sqsClient, cfg.SQS.AuditURL, auditSvc.Process, log, &wg, quit)

	// Start Workers
	validatorWorker.Start(5)
	enrichmentWorker.Start(5)
	auditWorker.Start(5)

	log.Info("All workers started")

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down workers...")
	close(quit)
	wg.Wait()

	log.Info("Workers exited gracefully")
}
