package query

import (
	"context"
)

type CreateVerifyEmailParams struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*VerifyEmail, error) {
	query := `INSERT INTO verify_emails (username, email, secret_code) VALUES ($1, $2, $3) 
            RETURNING id, username, email, secret_code, is_used, created_at, expired_at`

	row := q.db.QueryRowContext(ctx, query, arg.Username, arg.Email, arg.SecretCode)

	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

type UpdateVerifyEmailParams struct {
	SecretCode string `json:"secret_code"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (*VerifyEmail, error) {
	query := `
	   UPDATE verify_emails
	   SET
	       is_used = TRUE
	   WHERE
	       id = $1
	       AND secret_code = $2
	       AND is_used = FALSE
	       AND expired_at > now()
	   RETURNING id, username, email, secret_code, is_used, created_at, expired_at
	   `

	row := q.db.QueryRowContext(ctx, query, arg.ID, arg.SecretCode)

	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
