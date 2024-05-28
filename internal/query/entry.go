package query

import "context"

type CreateEntryParams struct {
	Account_ID int64 `json:"account_id"`
	Amount     int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (*Entry, error) {
	query := `INSERT INTO entries (account_id, amount) VALUES ($1, $2) 
            RETURNING id, account_id, amount, created_at`

	row := q.db.QueryRowContext(ctx, query, arg.Account_ID, arg.Amount)

	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Account_ID,
		&i.Amount,
		&i.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &i, nil

}

func (q *Queries) GetEntry(ctx context.Context, id int64) (*Entry, error) {
	query := "SELECT id, account_id, amount, created_at FROM entries WHERE id = $1"

	row := q.db.QueryRowContext(ctx, query, id)

	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Account_ID,
		&i.Amount,
		&i.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

type ListEntriesParams struct {
	Account_ID int64 `json:"account_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	query := `SELECT id, account_id, amount, created_at FROM entries WHERE account_id = $1
            ORDER BY id LIMIT $2 OFFSET $3`

	rows, err := q.db.QueryContext(ctx, query, arg.Account_ID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []Entry

	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.Account_ID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
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
