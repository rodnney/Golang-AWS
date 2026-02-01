package domain

import (
	"errors"
	"time"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusValidated TransactionStatus = "VALIDATED"
	StatusEnriched  TransactionStatus = "ENRICHED"
	StatusCompleted TransactionStatus = "COMPLETED"
	StatusFailed    TransactionStatus = "FAILED"
)

type TransactionType string

const (
	TypeDebit  TransactionType = "DEBIT"
	TypeCredit TransactionType = "CREDIT"
)

type Transaction struct {
	ID        string                 `json:"id"`
	AccountID string                 `json:"account_id"`
	Amount    float64                `json:"amount"`
	Currency  string                 `json:"currency"`
	Type      TransactionType        `json:"type"`
	Status    TransactionStatus      `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (t *Transaction) Validate() error {
	if t.ID == "" {
		return errors.New("id is required")
	}
	if t.AccountID == "" {
		return errors.New("account_id is required")
	}
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.Currency == "" {
		return ErrInvalidCurrency
	}
	if t.Type != TypeDebit && t.Type != TypeCredit {
		return errors.New("invalid transaction type")
	}
	return nil
}
