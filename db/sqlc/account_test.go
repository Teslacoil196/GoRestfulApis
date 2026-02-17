package db

import (
	"TeslaCoil196/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRamdonAccount(t *testing.T) Account {

	user := CreateRamdonUser(t)

	argument := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RamdonBalnce(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQuries.CreateAccount(context.Background(), argument)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, argument.Owner, account.Owner)
	require.Equal(t, argument.Balance, account.Balance)
	require.Equal(t, argument.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	CreateRamdonAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := CreateRamdonAccount(t)
	account2, err := testQuries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.Equal(t, account2.Balance, account1.Balance)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateRamdonAccount(t)
	bal := util.RamdonBalnce()
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: bal,
	}

	account2, err := testQuries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, account2.Balance, bal)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateRamdonAccount(t)
	err := testQuries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQuries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = CreateRamdonAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accoutns, err := testQuries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accoutns)

	for _, account := range accoutns {
		require.Equal(t, account.Owner, lastAccount.Owner)
	}
}
