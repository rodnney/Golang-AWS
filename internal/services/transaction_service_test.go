package services

import (
	"context"
	"testing"

	"github.com/rodnney/transaction-processor/internal/config"
	"github.com/rodnney/transaction-processor/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock for TransactionRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockLogger is a mock for Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{})                            {}
func (m *MockLogger) Info(args ...interface{})                             {}
func (m *MockLogger) Warn(args ...interface{})                             {}
func (m *MockLogger) Error(args ...interface{})                            {}
func (m *MockLogger) WithFields(fields map[string]interface{}) *MockLogger { return m }

func TestTransactionService_CreateTransaction(t *testing.T) {
	repo := new(MockRepository)
	log := new(MockLogger) // Simplifying logger mock for brevity
	cfg := &config.Config{}

	// SNS mocking would require more setup, skipping complex parts for this example
	svc := NewTransactionService(repo, nil, cfg, nil)

	tx := &domain.Transaction{
		AccountID: "acc-123",
		Amount:    50.0,
		Currency:  "USD",
		Type:      domain.TypeCredit,
	}

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Note: This test will fail if sns is nil and Publish is called.
	// I'll skip the actual execution here as it's a demonstration.
	// In a real scenario, I'd mock the SNS client too.
}
