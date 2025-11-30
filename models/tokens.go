package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"
)

const validationTokenLength = 32
const validationTokenTTL = 30 * time.Minute

type ValidationToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UsedAt    sql.NullTime
}

type ValidationTokenModel struct {
	DB *DB
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b), nil
}

func (m *ValidationTokenModel) Create(userID int64) (string, error) {
	token, err := generateRandomToken(validationTokenLength)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(validationTokenTTL)

	_, err = m.DB.Exec(
		`INSERT INTO validation_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID, token, expiresAt,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

var ErrInvalidOrExpiredToken = errors.New("invalid or expired token")

func (m *ValidationTokenModel) Consume(token string) (int64, error) {
	var (
		userID    int64
		expiresAt time.Time
		usedAt    sql.NullTime
	)

	err := m.DB.QueryRow(
		`SELECT user_id, expires_at, used_at
         FROM validation_tokens
         WHERE token = ?`,
		token,
	).Scan(&userID, &expiresAt, &usedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInvalidOrExpiredToken
		}
		return 0, err
	}

	if usedAt.Valid || time.Now().After(expiresAt) {
		return 0, ErrInvalidOrExpiredToken
	}

	_, err = m.DB.Exec(
		`UPDATE validation_tokens SET used_at = CURRENT_TIMESTAMP WHERE token = ?`,
		token,
	)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
