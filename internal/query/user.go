package query

import (
	"context"
	"database/sql"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	query := `INSERT INTO users (username, hashed_password, fullname, email) VALUES ($1, $2, $3, $4) 
                RETURNING username, hashed_password, fullname, email, password_changed_at, created_at`

	row := q.db.QueryRowContext(
		ctx,
		query,
		arg.Username,
		arg.HashedPassword,
		arg.Fullname,
		arg.Email,
	)

	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (q *Queries) GetUser(ctx context.Context, username string) (*User, error) {
	query := `SELECT username, hashed_password, fullname, email, password_changed_at, created_at 
                FROM users WHERE username = $1 LIMIT 1`

	row := q.db.QueryRowContext(ctx, query, username)

	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

type UpdateUserParams struct {
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	Username          string         `json:"username"`
	Fullname          sql.NullString `json:"fullname"`
	HashedPassword    sql.NullString `json:"hashed_password"`
	Email             sql.NullString `json:"email"`
	IsEmailVerified   sql.NullBool   `json:"is_email_verified"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (*User, error) {
	query := `UPDATE users SET
                fullname = COALESCE($1, fullname),
                hashed_password = COALESCE($2, hashed_password),
                email = COALESCE($3, email),
                password_changed_at = COALESCE($4, password_changed_at),
                is_email_verified = COALESCE($6, is_email_verified)
                WHERE username = $5
                RETURNING username, fullname, hashed_password, email, password_changed_at, created_at, is_email_verified`

	row := q.db.QueryRowContext(ctx, query, arg.Fullname, arg.HashedPassword, arg.Email, arg.PasswordChangedAt, arg.Username, arg.IsEmailVerified)

	var i User
	err := row.Scan(&i.Username, &i.Fullname, &i.HashedPassword, &i.Email, &i.PasswordChangedAt, &i.CreatedAt, &i.IsEmailVerified)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
