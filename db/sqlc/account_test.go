package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vnscriptkid/simple-bank-golang/util"
)

func createFakeAccount(t *testing.T) Account {
	params := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomInt(0, 1000),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), params)

	assert.NoError(t, err)
	assert.NotEmpty(t, acc)

	assert.Equal(t, params.Owner, acc.Owner)
	assert.Equal(t, params.Balance, acc.Balance)
	assert.Equal(t, params.Currency, acc.Currency)

	assert.NotEmpty(t, acc.CreatedAt)
	assert.NotEmpty(t, acc.ID)

	return acc
}

func TestCreateAccount(t *testing.T) {
	params := CreateAccountParams{
		Owner:    "Example Owner",
		Balance:  99,
		Currency: "USD",
	}

	acc, err := testQueries.CreateAccount(context.Background(), params)

	assert.NoError(t, err)
	assert.NotEmpty(t, acc)

	assert.Equal(t, params.Owner, acc.Owner)
	assert.Equal(t, params.Balance, acc.Balance)
	assert.Equal(t, params.Currency, acc.Currency)

	assert.NotEmpty(t, acc.CreatedAt)
	assert.NotEmpty(t, acc.ID)
}

func TestGetAccount(t *testing.T) {
	acc := createFakeAccount(t)

	gotAcc, err := testQueries.GetAccount(context.Background(), acc.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, acc)

	assert.Equal(t, acc.ID, gotAcc.ID)
	assert.Equal(t, acc.Balance, gotAcc.Balance)
	assert.Equal(t, acc.Currency, gotAcc.Currency)
	assert.Equal(t, acc.Owner, gotAcc.Owner)
	assert.WithinDuration(t, acc.CreatedAt, gotAcc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := createFakeAccount(t)

	params := UpdateAccountParams{
		ID:      acc.ID,
		Balance: acc.Balance + 1,
	}

	updatedAcc, err := testQueries.UpdateAccount(context.Background(), params)

	assert.NoError(t, err)
	assert.NotEmpty(t, updatedAcc)

	assert.Equal(t, acc.ID, updatedAcc.ID)
	assert.NotEqual(t, acc.Balance, updatedAcc.Balance)
	assert.Equal(t, acc.Currency, updatedAcc.Currency)
	assert.Equal(t, acc.Owner, updatedAcc.Owner)
	assert.Equal(t, acc.CreatedAt, updatedAcc.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	acc := createFakeAccount(t)

	err1 := testQueries.DeleteAccount(context.Background(), acc.ID)

	assert.NoError(t, err1)

	foundAcc, err2 := testQueries.GetAccount(context.Background(), acc.ID)

	assert.Error(t, err2)
	assert.EqualError(t, err2, sql.ErrNoRows.Error())
	assert.Empty(t, foundAcc)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFakeAccount(t)
	}

	params := ListAccountsParams{
		Offset: 5,
		Limit:  5,
	}

	accList, err := testQueries.ListAccounts(context.Background(), params)

	assert.NoError(t, err)
	assert.Len(t, accList, 5)

	for _, acc := range accList {
		assert.NotEmpty(t, acc)
	}
}
