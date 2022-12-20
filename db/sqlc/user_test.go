package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vnscriptkid/simple-bank-golang/util"
)

func createFakeUser(t *testing.T) User {
	params := CreateUserParams{
		// Owner:    util.RandomOwner(),
		// Balance:  util.RandomInt(0, 1000),
		// Currency: util.RandomCurrency(),
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, params.Username, user.Username)
	assert.Equal(t, params.HashedPassword, user.HashedPassword)
	assert.Equal(t, params.FullName, user.FullName)
	assert.Equal(t, params.Email, user.Email)

	assert.NotEmpty(t, user.CreatedAt)
	assert.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	user := createFakeUser(t)

	assert.NotEmpty(t, user)
}

func TestGetUser(t *testing.T) {
	user := createFakeUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), user.Username)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotUser)

	assert.Equal(t, user.CreatedAt, gotUser.CreatedAt)
	assert.Equal(t, user.Email, gotUser.Email)
	assert.Equal(t, user.FullName, gotUser.FullName)
	assert.Equal(t, user.HashedPassword, gotUser.HashedPassword)
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
}
