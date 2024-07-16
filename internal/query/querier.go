package query

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (*Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (*Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (*Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (*Account, error)
	GetEntry(ctx context.Context, id int64) (*Entry, error)
	GetTransfer(ctx context.Context, id int64) (*Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (*Account, error)
	GetUser(ctx context.Context, name string) (*User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (*User, error)
}
