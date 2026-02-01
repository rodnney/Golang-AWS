package services

import (
	"context"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type ValidatorService struct {
	repo   repository.TransactionRepository
	sns    *aws.SNSClient
	cfg    *config.Config
	logger logger.Logger
}

func NewValidatorService(repo repository.TransactionRepository, sns *aws.SNSClient, cfg *config.Config, logger logger.Logger) *ValidatorService {
	return &ValidatorService{
		repo:   repo,
		sns:    sns,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *ValidatorService) Process(ctx context.Context, tx *domain.Transaction) error {
	s.logger.Info("Validating transaction", "id", tx.ID)

	status := domain.StatusValidated
	if err := tx.Validate(); err != nil {
		status = domain.StatusFailed
		s.logger.Error("Transaction validation failed", "id", tx.ID, "error", err)
	}

	if err := s.repo.UpdateStatus(ctx, tx.ID, string(status)); err != nil {
		return err
	}

	if status == domain.StatusValidated {
		// In a real scenario, we might publish to another topic or just rely on the same topic and filter by status
		// For this simulation, we'll just log success
		s.logger.Info("Transaction validated successfully", "id", tx.ID)
	}

	return nil
}
