package query

import "context"

type CreateTransferParams struct {
	From_Account_ID int64 `json:"from_account_id"`
	To_Account_ID   int64 `json:"to_account_id"`
	Amount          int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (*Transfer, error) {
	query := `INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3) 
            RETURNING id, from_account_id, to_account_id, amount, created_at`

	row := q.db.QueryRowContext(ctx, query, arg.From_Account_ID, arg.To_Account_ID, arg.Amount)

	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.From_Account_ID,
		&i.To_Account_ID,
		&i.Amount,
		&i.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (q *Queries) GetTransfer(ctx context.Context, id int64) (*Transfer, error) {
	query := "SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers WHERE id = $1 LIMIT 1"

	row := q.db.QueryRowContext(ctx, query, id)

	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.From_Account_ID,
		&i.To_Account_ID,
		&i.Amount,
		&i.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

type ListTransfersParams struct {
	From_Account_ID int64 `json:"from_account_id"`
	To_Account_ID   int64 `json:"to_account_id"`
	Limit           int32 `json:"limit"`
	Offset          int32 `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
                WHERE from_account_id = $1 OR to_account_id = $2
                ORDER BY id LIMIT $3 OFFSET $4`

	rows, err := q.db.QueryContext(
		ctx,
		query,
		arg.From_Account_ID,
		arg.To_Account_ID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Transfer
	for rows.Next() {
		var i Transfer
		err := rows.Scan(
			&i.ID,
			&i.From_Account_ID,
			&i.To_Account_ID,
			&i.Amount,
			&i.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
