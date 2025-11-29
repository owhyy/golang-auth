package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail     = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
	IsValid      bool
	CreatedAt    string
}

type UserModel struct {
	DB *DB
}

func (m *UserModel) Create(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, string(hashedPassword),
	)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) error {
	var passwordHash string

	err := m.DB.QueryRow(
		"SELECT password_hash FROM users WHERE email = ?",
		email,
	).Scan(&passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrInvalidCredentials
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}

func (m *UserModel) EmailExists(email string) (bool, error) {
	var exists bool
	err := m.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)

	return exists, err
}
