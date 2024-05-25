package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/earnstein/GO/greenlight/internal/validator"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// MOVIE STRUCTS

type (
	Movie struct {
		ID        int64     `json:"id"`
		Title     string    `json:"title"`
		Year      int32     `json:"year,omitempty"`
		Runtime   Runtime   `json:"runtime,omitempty"`
		Genres    []string  `json:"genres,omitempty"`
		CreatedAt time.Time `json:"-"`
		Version   int32     `json:"version"`
	}

	MovieModel struct {
		DB *sql.DB
	}

	Models struct {
		Movies MovieModel
	}
)

func NewModels(db *sql.DB) *Models {
	return &Models{
		Movies: MovieModel{DB: db},
	}
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

func (m *MovieModel) Insert(movie *Movie) error {
	stmt := `INSERT INTO movies (title, year, runtime, genres)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	row := m.DB.QueryRowContext(ctx, stmt, args...)
	err := row.Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	return err
}

func (m *MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}
	stmt := `SELECT id, title, year, runtime, genres, created_at, version
			FROM movies
			WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.CreatedAt,
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m *MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {

	stmt := fmt.Sprintf(`
		SELECT id, title, year, runtime, genres, created_at, version
		FROM movies
		WHERE (to_tsvector('english', title) @@ plainto_tsquery('english', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, stmt, title, pq.Array(genres))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.CreatedAt,
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
func (m *MovieModel) Update(movie *Movie) error {

	stmt := `UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	stmt := ` DELETE FROM movies
			WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
