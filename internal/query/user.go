package query

import "context"

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
