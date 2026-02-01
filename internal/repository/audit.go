package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/domain"
)

type AuditRepository interface {
	Save(ctx context.Context, audit *domain.AuditLog) error
}

type S3AuditRepository struct {
	client *aws.S3Client
	cfg    *config.Config
}

func NewS3AuditRepository(client *aws.S3Client, cfg *config.Config) *S3AuditRepository {
	return &S3AuditRepository{
		client: client,
		cfg:    cfg,
	}
}

func (r *S3AuditRepository) Save(ctx context.Context, audit *domain.AuditLog) error {
	data, err := json.Marshal(audit)
	if err != nil {
		return fmt.Errorf("failed to marshal audit log: %w", err)
	}

	// Key format: logs/YYYY/MM/DD/transaction_ID.json
	now := time.Now()
	key := fmt.Sprintf("logs/%d/%02d/%02d/%s_%s.json",
		now.Year(), now.Month(), now.Day(), audit.TransactionID, audit.ID)

	return r.client.Upload(ctx, r.cfg.S3.Bucket, key, data)
}
