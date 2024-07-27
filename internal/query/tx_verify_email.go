package query

import (
	"context"
	"database/sql"
)

type VerifyEmailTxParams struct {
	SecretCode string
	EmailID    int64
}

type VerifyEmailTxResult struct {
	User        *User
	VerifyEmail *VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (*VerifyEmailTxResult, error) {
	txResult := &VerifyEmailTxResult{}

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txResult.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		txResult.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: txResult.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		})

		return err
	})

	return txResult, err
}
