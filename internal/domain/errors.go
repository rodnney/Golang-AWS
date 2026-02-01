package domain

import "errors"

var (
	ErrInvalidTransaction  = errors.New("invalid transaction")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidAmount       = errors.New("invalid amount: must be greater than zero")
	ErrInvalidCurrency     = errors.New("invalid currency: cannot be empty")
)
