package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefrestToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIP     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error) {
	query := `INSERT INTO sessions (id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at`

	row := q.db.QueryRowContext(
		ctx,
		query,
		arg.ID,
		arg.Username,
		arg.RefrestToken,
		arg.UserAgent,
		arg.ClientIP,
		arg.IsBlocked,
		arg.ExpiresAt,
	)

	var i Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefrestToken,
		&i.UserAgent,
		&i.ClientIP,
		&i.IsBlocked,
		&i.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	query := "SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at FROM sessions WHERE id = $1"
	row := q.db.QueryRowContext(ctx, query, id)

	var i Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefrestToken,
		&i.UserAgent,
		&i.ClientIP,
		&i.IsBlocked,
		&i.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &i, nil
}
