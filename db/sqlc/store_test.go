package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(db)

	account1 := CreateRamdonAccount(t)
	account2 := CreateRamdonAccount(t)
	fmt.Println("Balance of accounts Before ->", account1.Balance, account2.Balance)

	numberOfConcurrentTransactions := 5
	amount := int64(5)
	results := make(chan TransferTxResult)
	errr := make(chan error)

	// very important that you send and read from the channels in right order
	// otherwise we create a deadlock
	for i := 0; i < numberOfConcurrentTransactions; i++ {
		txName := fmt.Sprintf(" tx %d", i+1)
		go func() {
			transfer := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TranferTx(ctx, transfer)

			errr <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < numberOfConcurrentTransactions; i++ {
		err := <-errr
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.NotEmpty(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.FromEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, ToEntry.AccountID, account1.ID)
		require.Equal(t, ToEntry.Amount, -amount)
		require.NotEmpty(t, ToEntry.CreatedAt)
		require.NotEmpty(t, ToEntry.ID)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// TODO update the user account and balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		differ1 := account1.Balance - fromAccount.Balance
		differ2 := toAccount.Balance - account2.Balance
		require.Equal(t, differ1, differ2)
		require.True(t, differ1 > 0)
		require.True(t, differ2%amount == 0)

		k := int(differ1 / amount)
		require.True(t, k >= 1 && k <= numberOfConcurrentTransactions)

		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQuries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQuries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println("Balance of accounts After ->", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-(amount*int64(numberOfConcurrentTransactions)), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(amount*int64(numberOfConcurrentTransactions)), updatedAccount2.Balance)

}
