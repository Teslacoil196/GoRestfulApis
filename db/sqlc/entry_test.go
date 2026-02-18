package db

import (
	"TeslaCoil196/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateAccountForEntries(t *testing.T) {
	account := CreateRamdonAccount(t)

	TestAccount3 = account
}

func CreateRandomEntry(t *testing.T) Entry {

	if !accountCreatedForEntry {
		CreateAccountForEntries(t)
		accountCreatedForEntry = true
	}

	agu := CreateEntryParams{
		AccountID: TestAccount3.ID,
		Amount:    util.RamdonBalnce(),
	}

	entry, err := testQuries.CreateEntry(context.Background(), agu)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, TestAccount3.ID)
	//require.Equal(t, 123456, TestAccount3.ID)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestDeleteEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	err := testQuries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry1, err := testQuries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry1)
}

func TestGetEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	entry1, err := testQuries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.Equal(t, entry.AccountID, entry1.AccountID)
	require.Equal(t, entry.ID, entry1.ID)
	require.Equal(t, entry.Amount, entry1.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry1.CreatedAt, time.Second)

}
