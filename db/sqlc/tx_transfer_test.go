package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTxn(t *testing.T) {
	sqlStore := NewStore(testDB)

	account1 := createFakeAccount(t)
	account2 := createFakeAccount(t)

	amount := int64(10)

	errors := make(chan error)
	results := make(chan TransferTxResult)

	n := 10

	for i := 0; i < n; i++ {
		go func() {
			result, err := sqlStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors

		assert.NoError(t, err)

		result := <-results

		assert.NotEmpty(t, result)

		// FromAccount Account  `json:"from_account"`
		// ToAccount   Account  `json:"to_account"`

		// Transfer    Transfer `json:"transfer"`
		transfer := result.Transfer
		assert.Equal(t, transfer.Amount, amount)
		assert.Equal(t, transfer.FromAccountID, account1.ID)
		assert.Equal(t, transfer.ToAccountID, account2.ID)
		assert.NotEmpty(t, transfer.ID)
		assert.NotEmpty(t, transfer.CreatedAt)

		// FromEntry   Entry    `json:"from_entry"`
		fromEntry := result.FromEntry
		assert.Equal(t, fromEntry.Amount, -amount)
		assert.Equal(t, fromEntry.AccountID, account1.ID)
		assert.NotEmpty(t, fromEntry.ID)
		assert.NotEmpty(t, fromEntry.CreatedAt)

		// ToEntry     Entry    `json:"to_entry"`
		toEntry := result.ToEntry
		assert.Equal(t, toEntry.Amount, amount)
		assert.Equal(t, toEntry.AccountID, account2.ID)
		assert.NotEmpty(t, toEntry.ID)
		assert.NotEmpty(t, toEntry.CreatedAt)
	}
}
