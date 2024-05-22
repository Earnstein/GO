package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/earnstein/GO/greenlight/internal/validator"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
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

func NewModels(db *sql.DB) Models {
	return Models{
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

func (m MovieModel) Insert(movie *Movie) error {
	stmt := `INSERT INTO movies (title, year, runtime, genres)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	row := m.DB.QueryRow(stmt, args...)
	err := row.Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	return err
}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}
	stmt := `SELECT id, title, year, runtime, genres, created_at, version
			FROM movies
			WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(stmt, id).Scan(
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

func (m MovieModel) Update(movie *Movie) error {

	stmt := `UPDATE movies
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}
	err := m.DB.QueryRow(stmt, args...).Scan(&movie.Version)
	return err
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	stmt := ` DELETE FROM movies
			WHERE id = $1`

	result, err := m.DB.Exec(stmt, id)
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
