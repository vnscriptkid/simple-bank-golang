package db

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTxn(t *testing.T) {
	sqlStore := NewStore(testDB)

	account1 := createFakeAccount(t)
	account2 := createFakeAccount(t)

	fmt.Println(">> Before ", account1.Balance, account2.Balance)

	amount := int64(100)

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

		// FromAccount Account  `json:"from_account"`
		fromAcc := result.FromAccount
		assert.NotEmpty(t, fromAcc)
		assert.Equal(t, fromAcc.ID, account1.ID)

		// ToAccount   Account  `json:"to_account"`
		toAcc := result.ToAccount
		assert.NotEmpty(t, toAcc)
		assert.Equal(t, toAcc.ID, account2.ID)

		diff1 := account1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - account2.Balance
		assert.Equal(t, diff1, diff2)
		assert.Equal(t, diff1, int64(i+1)*amount)
		fmt.Println(">>>> During "+strconv.Itoa(i)+": ", fromAcc.Balance, toAcc.Balance)
	}

	account1Final, _ := sqlStore.GetAccount(context.Background(), account1.ID)
	account2Final, _ := sqlStore.GetAccount(context.Background(), account2.ID)

	fmt.Println(">> After ", account1Final.Balance, account2Final.Balance)

	assert.Equal(t, account1Final.Balance, account1.Balance-int64(n)*amount)
	assert.Equal(t, account2Final.Balance, account2.Balance+int64(n)*amount)
}

func TestDeadlockInTransferTxn(t *testing.T) {
	sqlStore := NewStore(testDB)

	account1 := createFakeAccount(t)
	account2 := createFakeAccount(t)

	fmt.Println(">> Before ", account1.Balance, account2.Balance)

	amount := int64(100)

	errors := make(chan error)

	n := 10

	for i := 0; i < n; i++ {
		fromAccId := account1.ID
		toAccId := account2.ID

		if i%2 == 1 {
			fromAccId = account2.ID
			toAccId = account1.ID
		}

		go func() {

			_, err := sqlStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccId,
				ToAccountID:   toAccId,
				Amount:        amount,
			})

			errors <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors

		assert.NoError(t, err)
	}

	account1Final, _ := sqlStore.GetAccount(context.Background(), account1.ID)
	account2Final, _ := sqlStore.GetAccount(context.Background(), account2.ID)

	fmt.Println(">> After ", account1Final.Balance, account2Final.Balance)

	assert.Equal(t, account1Final.Balance, account1.Balance)
	assert.Equal(t, account2Final.Balance, account2.Balance)
}
