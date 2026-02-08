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

	numberOfConcurrentTransactions := 5
	amount := int64(5)
	fmt.Print("1 ======================================")
	results := make(chan TransferTxResult)
	errr := make(chan error)

	// very important that you send and read from the channels in right order
	// otherwise we create a deadlock
	for i := 0; i < numberOfConcurrentTransactions; i++ {
		go func() {
			transfer := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			result, err := store.TranferTx(context.Background(), transfer)

			errr <- err
			results <- result
		}()
	}

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

	}

}
