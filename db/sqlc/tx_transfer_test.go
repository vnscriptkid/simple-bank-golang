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

	n := 1

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

		// // FromAccount Account  `json:"from_account"`
		// fromAcc := result.FromAccount
		// assert.Equal(t, account1.ID, fromAcc.ID)

		// // ToAccount   Account  `json:"to_account"`
		// toAcc := result.FromAccount
		// assert.Equal(t, account1.ID, toAcc.ID)

		// // Transfer    Transfer `json:"transfer"`
		// transfer := result.Transfer
		// assert.Equal(t, transfer.Amount, amount)
		// assert.Equal(t, transfer.FromAccountID, account1.ID)
		// assert.Equal(t, transfer.ToAccountID, account2.ID)
		// assert.Greater(t, transfer.ID, 0)
		// assert.NotEmpty(t, transfer.CreatedAt)

		// // FromEntry   Entry    `json:"from_entry"`
		// fromEntry := result.FromEntry
		// assert.Equal(t, fromEntry.Amount, -amount)
		// assert.Equal(t, fromEntry.AccountID, account1.ID)
		// assert.Greater(t, fromEntry.ID, 0)
		// assert.NotEmpty(t, fromEntry.CreatedAt)

		// // ToEntry     Entry    `json:"to_entry"`
		// toEntry := result.ToEntry
		// assert.Equal(t, toEntry.Amount, amount)
		// assert.Equal(t, toEntry.AccountID, account2.ID)
		// assert.Greater(t, toEntry.ID, 0)
		// assert.NotEmpty(t, toEntry.CreatedAt)
	}
}
