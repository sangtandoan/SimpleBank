package query

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(context.Context, TransferTxParams) (*TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx) // creates a new Queries with tx so that all queries are running in that tx
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	From_Account_ID int64 `json:"from_account_id"`
	To_Account_ID   int64 `json:"to_account_id"`
	Amount          int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer     *Transfer `json:"transfer"`
	From_Account *Account  `json:"from_account"`
	To_Account   *Account  `json:"to_account"`
	From_Entry   *Entry    `json:"from_entry"`
	To_Entry     *Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) TransferTx(
	ctx context.Context,
	arg TransferTxParams,
) (*TransferTxResult, error) {
	// remembers that pointer needs to point to some address, need init for pointer
	var result *TransferTxResult
	result = &TransferTxResult{}

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			From_Account_ID: arg.From_Account_ID,
			To_Account_ID:   arg.To_Account_ID,
			Amount:          arg.Amount,
		})
		if err != nil {
			return err
		}

		result.From_Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			Account_ID: arg.From_Account_ID,
			Amount:     -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.To_Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			Account_ID: arg.To_Account_ID,
			Amount:     arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		if arg.From_Account_ID < arg.To_Account_ID {
			result.From_Account, result.To_Account, err = addMoney(
				ctx,
				q,
				arg.From_Account_ID,
				-arg.Amount,
				arg.To_Account_ID,
				arg.Amount,
			)
		} else {
			result.To_Account, result.From_Account, err = addMoney(ctx, q, arg.To_Account_ID, arg.Amount, arg.From_Account_ID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 *Account, account2 *Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
