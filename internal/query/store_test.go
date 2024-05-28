package query

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan *TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				From_Account_ID: account1.ID,
				To_Account_ID:   account2.ID,
				Amount:          amount,
			})

			errs <- err
			results <- result
		}()
	}
	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.From_Account_ID)
		require.Equal(t, account2.ID, transfer.To_Account_ID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from_entry
		from_entry := result.From_Entry
		require.NotEmpty(t, from_entry)
		require.Equal(t, account1.ID, from_entry.Account_ID)
		require.Equal(t, -amount, (from_entry).Amount)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), from_entry.ID)
		require.NoError(t, err)

		// check to_entry
		to_entry := result.To_Entry
		require.NotEmpty(t, to_entry)
		require.Equal(t, account2.ID, to_entry.Account_ID)
		require.Equal(t, amount, to_entry.Amount)
		require.NotZero(t, to_entry.ID)
		require.NotZero(t, to_entry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), to_entry.ID)
		require.NoError(t, err)

		// check accounts' balance
		fromAccount := result.From_Account
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.To_Account
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-amount*int64(n), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+amount*int64(n), updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i >= n/2 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				From_Account_ID: fromAccountID,
				To_Account_ID:   toAccountID,
				Amount:          amount,
			})

			errs <- err
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
