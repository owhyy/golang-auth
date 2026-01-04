package main

import (
	"html/template"
	"io/fs"
	"owhyy/simple-auth/internal/models"
	"owhyy/simple-auth/ui"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

type paginationData struct {
	CurrentPage int
	PerPage     int
	TotalPages  int
	Prev        int
	Next        int
}

type templateData struct {
	User            models.User
	Error           string
	Token           string
	IsAuthenticated bool
	BaseURL         string
	Posts           []models.Post
	Post            models.Post
	Pagination      paginationData
}
