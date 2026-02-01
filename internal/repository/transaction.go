package repository

import (
	"context"

	"github.com/rodnney/transaction-processor/internal/aws"
	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	// Note: List is mentioned in prompt but we'll focus on the core flow first
	// List(ctx context.Context, limit int) ([]*domain.Transaction, error)
}

type DynamoDBTransactionRepository struct {
	client *aws.DynamoDBClient
	cfg    *config.Config
}

func NewDynamoDBTransactionRepository(client *aws.DynamoDBClient, cfg *config.Config) *DynamoDBTransactionRepository {
	return &DynamoDBTransactionRepository{
		client: client,
		cfg:    cfg,
	}
}

func (r *DynamoDBTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	return r.client.PutItem(ctx, r.cfg.DynamoDB.Table, tx)
}

func (r *DynamoDBTransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.client.GetItem(ctx, r.cfg.DynamoDB.Table, id, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *DynamoDBTransactionRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	return r.client.UpdateItem(ctx, r.cfg.DynamoDB.Table, id, updates)
}
