package query

import (
	"context"
	"testing"
	"time"

	"github.com/FrostJ143/simplebank/internal/query/utils"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, from_account *Account, to_account *Account) *Transfer {
	arg := CreateTransferParams{
		From_Account_ID: from_account.ID,
		To_Account_ID:   to_account.ID,
		Amount:          utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.From_Account_ID, transfer.From_Account_ID)
	require.Equal(t, arg.To_Account_ID, transfer.To_Account_ID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.From_Account_ID, transfer2.From_Account_ID)
	require.Equal(t, transfer1.To_Account_ID, transfer2.To_Account_ID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
		createRandomTransfer(t, account2, account1)
	}

	arg := ListTransfersParams{
		From_Account_ID: account1.ID,
		To_Account_ID:   account1.ID,
		Limit:           5,
		Offset:          5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(
			t,
			transfer.From_Account_ID == arg.From_Account_ID ||
				transfer.To_Account_ID == arg.To_Account_ID,
		)
	}
}
