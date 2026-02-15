package db

import (
	"TeslaCoil196/util"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRamdonUser(t *testing.T) User {

	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	argument := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQuries.CreateUser(context.Background(), argument)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, argument.Username, user.Username)
	require.Equal(t, argument.HashedPassword, user.HashedPassword)
	require.Equal(t, argument.FullName, user.FullName)
	require.Equal(t, argument.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.True(t, user.LastPassReset.IsZero())

	return user

}

func TestCreateUser(t *testing.T) {
	CreateRamdonUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRamdonUser(t)
	user2, err := testQuries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.HashedPassword, user1.HashedPassword)
	require.Equal(t, user2.FullName, user1.FullName)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.LastPassReset, user2.LastPassReset, time.Second)

}
