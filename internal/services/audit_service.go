package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type AuditService struct {
	repo   repository.AuditRepository
	logger logger.Logger
}

func NewAuditService(repo repository.AuditRepository, logger logger.Logger) *AuditService {
	return &AuditService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AuditService) Process(ctx context.Context, tx *domain.Transaction) error {
	s.logger.Info("Auditing transaction", "id", tx.ID)

	audit := &domain.AuditLog{
		ID:            uuid.New().String(),
		TransactionID: tx.ID,
		Action:        string(tx.Status),
		Status:        "SUCCESS",
		Message:       "Transaction processed step: " + string(tx.Status),
		Timestamp:     time.Now(),
	}

	err := s.repo.Save(ctx, audit)
	if err != nil {
		s.logger.Error("Failed to save audit log", "id", tx.ID, "error", err)
		return err
	}

	s.logger.Info("Transaction audited successfully", "id", tx.ID)
	return nil
}
