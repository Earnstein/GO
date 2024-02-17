// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: entry.sql

package db

import (
	"context"
)

const getEntry = `-- name: GetEntry :one
SELECT id, amount, account_id, created_at FROM entries
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.AccountID,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, amount, account_id, created_at FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.AccountID,
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

const createEntry = `-- name: createEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1 , $2 
) RETURNING id, amount, account_id, created_at
`

type createEntryParams struct {
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (q *Queries) createEntry(ctx context.Context, arg createEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.AccountID,
		&i.CreatedAt,
	)
	return i, err
}
