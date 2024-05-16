package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"hashedpassword"`
	Created        time.Time `json:"created"`
}

type UserRequestBody struct {
	Name     string `json:"name" validate:"required" binding:"required"`
	Email    string `json:"email" validate:"required" binding:"required"`
	Password string `json:"password" validate:"required" binding:"required"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES
		(?, ?, ?, UTC_TIMESTAMP())
		`
	_, err = m.DB.Exec(stmt, name, email, string(hashed_password))

	if err != nil {
		var queryErr *mysql.MySQLError
		if errors.As(err, &queryErr) {
			if queryErr.Number == 1062 && strings.Contains(queryErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exist(id int) (bool, error) {
	return false, nil
}
