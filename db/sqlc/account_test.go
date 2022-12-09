package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
