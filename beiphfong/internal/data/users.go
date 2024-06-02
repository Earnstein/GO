package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/earnstein/GO/greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

func NewUser(name, email string, activated bool) *User {
	return &User{
		Username:  name,
		Email:     email,
		Activated: activated,
	}
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	password = strings.Trim(password, " ")
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
	v.Check(strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ"), "password", "password should contain at least one uppercase letter")
	v.Check(strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz"), "password", "password should contain at least one lowercase letter")
	v.Check(strings.ContainsAny(password, "0123456789"), "password", "password should contain at least one number")
	v.Check(validator.Matches(password, validator.SpecialCharRegex), "password", "password should contain at least one special character")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "name", "must be provided")
	v.Check(len(user.Username) <= 250, "name", "must not be more than 250 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	stmt := `INSERT INTO users (username, email, password_hash, activated)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`
	args := []interface{}{user.Username, user.Email, user.Password.hash, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	stmt := `SELECT id, created_at, username, email, password_hash, activated, version
			FROM users
			WHERE email = $1`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	stmt := `UPDATE users
			SET username = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version`
	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
