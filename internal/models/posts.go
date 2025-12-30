package models

import (
	"time"
)

type PostStatus string

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
	Status        string
	PublishedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FeaturedImage string
}
