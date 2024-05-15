package models

import (
	"database/sql"
	"errors"
)

type Snippet struct {
	ID      *int    `json:"id"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Created *string `json:"created"`
	Expires *string `json:"expires"`
}

type SnippetModel struct {
	DB *sql.DB
}

type SnippetBody struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required, min=10, max=100"`
	Expires string `json:"expires" validate:"required"`
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES
			(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created , expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	data := &Snippet{}

	if err := row.Scan(
		&data.ID,
		&data.Title,
		&data.Content,
		&data.Created,
		&data.Expires,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return data, nil
}

func (m *SnippetModel) GetLatest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		data := &Snippet{}
		if err = rows.Scan(
			&data.ID,
			&data.Title,
			&data.Content,
			&data.Created,
			&data.Expires,
		); err != nil {
			return nil, err
		}

		snippets = append(snippets, data)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
