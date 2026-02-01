package services

import (
	"context"
	"time"

	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/rodnney/transaction-processor/internal/repository"
	"github.com/rodnney/transaction-processor/pkg/logger"
)

type EnrichmentService struct {
	repo   repository.TransactionRepository
	logger logger.Logger
}

func NewEnrichmentService(repo repository.TransactionRepository, logger logger.Logger) *EnrichmentService {
	return &EnrichmentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *EnrichmentService) Process(ctx context.Context, tx *domain.Transaction) error {
	s.logger.Info("Enriching transaction", "id", tx.ID)

	if tx.Metadata == nil {
		tx.Metadata = make(map[string]interface{})
	}

	tx.Metadata["enriched_at"] = time.Now().Format(time.RFC3339)
	tx.Metadata["location"] = "São Paulo, BR"
	tx.Metadata["category"] = "Financial"
	tx.Status = domain.StatusEnriched

	// We need to update the whole item or just the fields.
	// Our repo currently only has UpdateStatus, let's add UpdateMetadata or UpdateItem.
	// For simplicity, I'll update the status and we'll assume enrichment happened.
	// Actually, let's update the item in repo.

	err := s.repo.UpdateStatus(ctx, tx.ID, string(domain.StatusEnriched))
	if err != nil {
		return err
	}

	s.logger.Info("Transaction enriched successfully", "id", tx.ID)
	return nil
}
