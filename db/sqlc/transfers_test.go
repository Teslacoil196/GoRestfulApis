package db

import (
	"TeslaCoil196/util"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// var (
// 	TestAccount1
// 	TestAccount2
// )

//var accountsCreated = false

func CreateAccountsForTransferTests(t *testing.T) {
	account1 := CreateRamdonAccount(t)
	account2 := CreateRamdonAccount(t)

	TestAccount1 = account1
	TestAccount2 = account2
}

func CreateRandomTransfer(t *testing.T) Transfer {

	if !accountsCreatedForTrasnfers {
		CreateAccountsForTransferTests(t)
		accountsCreatedForTrasnfers = true
	}

	arug := CreateTransferParams{
		FromAccountID: TestAccount1.ID,
		ToAccountID:   TestAccount2.ID,
		Amount:        util.RamdonBalnce(),
	}

	tranfer, err := testQuries.CreateTransfer(context.Background(), arug)

	require.NoError(t, err)
	require.NotEmpty(t, tranfer)

	require.Equal(t, tranfer.FromAccountID, TestAccount1.ID)
	require.Equal(t, tranfer.ToAccountID, TestAccount2.ID)

	require.NotZero(t, tranfer.ID)
	require.NotZero(t, tranfer.CreatedAt)

	return tranfer
}

func TestCreateTransfer(t *testing.T) {
	CreateRandomTransfer(t)
}

func TestGetTranfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	transfer1, err := testQuries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer1)

	require.Equal(t, transfer.ID, transfer1.ID)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer.Amount, transfer1.Amount)
	require.WithinDuration(t, transfer.CreatedAt, transfer1.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)
	err := testQuries.DeleteTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)

	transfer1, err := testQuries.GetTransfer(context.Background(), transfer.ID)

	require.Empty(t, transfer1)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	aggo := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}

	transfers, err := testQuries.ListTransfers(context.Background(), aggo)

	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}

func TestListTransfersFromAccount(t *testing.T) {
	trans := []Transfer{}

	fmt.Print(trans)

	for i := 0; i < 10; i++ {
		t := CreateRandomTransfer(t)
		trans = append(trans, t)
	}

	fmt.Print(trans)

	aggo := ListTransfersFromAccountParams{
		FromAccountID: TestAccount1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQuries.ListTransfersFromAccount(context.Background(), aggo)

	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}

func TestListTransfersTOAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	aggo := ListTransfersToAccountParams{
		ToAccountID: TestAccount2.ID,
		Limit:       5,
		Offset:      5,
	}

	transfers, err := testQuries.ListTransfersToAccount(context.Background(), aggo)

	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}

func TestListTransfersFromAccountTOAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t)
	}

	aggo := ListTransfersFromAccountToAccountParams{
		FromAccountID: TestAccount1.ID,
		ToAccountID:   TestAccount2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQuries.ListTransfersFromAccountToAccount(context.Background(), aggo)

	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}

func TestUpdateTransfer(t *testing.T) {
	transfer := CreateRandomTransfer(t)

	bal := util.RamdonBalnce()

	arg := UpdateTransferParams{
		ID:     transfer.ID,
		Amount: bal,
	}

	transfer1, err := testQuries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, transfer.ID, transfer1.ID)
	require.Equal(t, transfer.FromAccountID, transfer1.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer1.ToAccountID)
	require.Equal(t, transfer1.Amount, bal)

}
