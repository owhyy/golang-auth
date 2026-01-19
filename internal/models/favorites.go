package models

import (
	"errors"
	"time"
)

var (
	ErrAlreadyFavorited = errors.New("post already favorited by user")
	ErrNotFavorited     = errors.New("post not favorited by user")
)

type Favorite struct {
	ID        uint
	UserID    uint
	PostID    uint
	CreatedAt time.Time
}

type FavoriteModel struct {
	DB *DB
}

func (m *FavoriteModel) Add(userID, postID uint) error {
	query := `
		INSERT INTO favorites (user_id, post_id)
		VALUES (?, ?)
	`

	result, err := m.DB.Exec(query, userID, postID)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: favorites.user_id, favorites.post_id" {
			return ErrAlreadyFavorited
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrAlreadyFavorited
	}

	return nil
}

func (m *FavoriteModel) Remove(userID, postID uint) error {
	query := `
		DELETE FROM favorites 
		WHERE user_id = ? AND post_id = ?
	`

	result, err := m.DB.Exec(query, userID, postID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFavorited
	}

	return nil
}

func (m *FavoriteModel) IsFavorited(userID, postID uint) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM favorites
		WHERE user_id = ? AND post_id = ?
	`

	var count int
	err := m.DB.QueryRow(query, userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *FavoriteModel) GetFavoriteCount(postID uint) (int, error) {
	query := `
		SELECT COUNT(1)
		FROM favorites
		WHERE post_id = ?
	`

	var count int
	err := m.DB.QueryRow(query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
