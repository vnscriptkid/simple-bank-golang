package db

import (
	"context"
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
