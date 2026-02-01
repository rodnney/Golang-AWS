package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type TransactionService struct {
	repo   repository.TransactionRepository
	sns    *aws.SNSClient
	cfg    *config.Config
	logger logger.Logger
}

func NewTransactionService(repo repository.TransactionRepository, sns *aws.SNSClient, cfg *config.Config, logger logger.Logger) *TransactionService {
	return &TransactionService{
		repo:   repo,
		sns:    sns,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	tx.ID = uuid.New().String()
	tx.Status = domain.StatusPending
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	if err := tx.Validate(); err != nil {
		return err
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return err
	}

	s.logger.Info("Transaction created", "id", tx.ID)

	if err := s.sns.Publish(ctx, s.cfg.SNS.TopicARN, tx); err != nil {
		s.logger.Error("Failed to publish transaction event", "error", err)
		return err
	}

	return nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}
