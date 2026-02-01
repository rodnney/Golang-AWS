package domain

import "time"

type AuditLog struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Action        string    `json:"action"`
	Status        string    `json:"status"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
}
