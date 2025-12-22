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
	ID           int64
	Email        string
	PasswordHash string
	IsValid      bool
	CreatedAt    string
}

type UserModel struct {
	DB *DB
}

func (m *UserModel) Create(email, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := m.DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, string(hashedPassword),
	)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return 0, ErrDuplicateEmail
		}
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (int64, error) {
	var id int64
	var passwordHash string

	err := m.DB.QueryRow(
		"SELECT id, password_hash FROM users WHERE email = ?",
		email,
	).Scan(&id, &passwordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return 0, ErrInvalidCredentials
	}

	return id, nil
}

func (m *UserModel) EmailExists(email string) (bool, error) {
	var exists bool
	err := m.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)

	return exists, err
}

func (m *UserModel) ValidateByID(id int64) error {
	_, err := m.DB.Exec(
		`UPDATE users SET is_valid = 1 WHERE id = ?`,
		id,
	)
	return err
}

func (m *UserModel) GetEmailByID(id int64) (string, error) {
	var email string
	err := m.DB.QueryRow(
		"SELECT email FROM users WHERE id = ?",
		id,
	).Scan(&email)

	return email, err
}

func (m *UserModel) GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := m.DB.QueryRow(
		"SELECT id, email, is_valid, created_at FROM users WHERE id = ?",
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.IsValid,
		&user.CreatedAt,
	)
	return user, err
}
