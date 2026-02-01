package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransaction_Validate(t *testing.T) {
	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				ID:        "123",
				AccountID: "acc-1",
				Amount:    100.0,
				Currency:  "BRL",
				Type:      TypeDebit,
			},
			wantErr: false,
		},
		{
			name: "missing id",
			tx: Transaction{
				AccountID: "acc-1",
				Amount:    100.0,
				Currency:  "BRL",
				Type:      TypeDebit,
			},
			wantErr: true,
		},
		{
			name: "invalid amount",
			tx: Transaction{
				ID:        "123",
				AccountID: "acc-1",
				Amount:    -10.0,
				Currency:  "BRL",
				Type:      TypeDebit,
			},
			wantErr: true,
		},
		{
			name: "missing currency",
			tx: Transaction{
				ID:        "123",
				AccountID: "acc-1",
				Amount:    100.0,
				Type:      TypeDebit,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
