package models

import (
	"time"
	"errors"
	"database/sql"
)

type PostStatus string

var ErrRecordNotFound = errors.New("record not found")

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

type Post struct {
	ID            uint
	Title         string
	Slug          string
	Content       string
	Excerpt       string
	AuthorID      uint
	Status        PostStatus
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FeaturedImage *string
}

type PostModel struct {
	DB *DB
}

func (m *PostModel) GetPublished(limit uint) ([]Post, error) {
	query := `
		SELECT
			id,
			title,
			slug,
			content,
			excerpt,
			author_id,
			status,
			published_at,
			created_at,
			updated_at,
			featured_image
		FROM posts
		WHERE status = 'published'
		ORDER BY published_at DESC
		LIMIT $1
	`

	rows, err := m.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post

		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.Excerpt,
			&p.AuthorID,
			&p.Status,
			&p.PublishedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.FeaturedImage,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *PostModel) GetBySlug(slug string) (*Post, error) {
	query := `
		SELECT
			id,
			title,
			slug,
			content,
			excerpt,
			author_id,
			status,
			published_at,
			created_at,
			updated_at,
			featured_image
		FROM posts
		WHERE slug = $1
		  AND status = 'published'
		LIMIT 1
	`

	var p Post

	err := m.DB.QueryRow(query, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.Excerpt,
		&p.AuthorID,
		&p.Status,
		&p.PublishedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.FeaturedImage,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &p, nil
}
