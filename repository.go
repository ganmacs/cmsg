package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

type Repository struct {
	Path, Name, ShortPath string
}

func newRepository(repo, root string) *Repository {
	return &Repository{
		Path:      repo,
		Name:      filepath.Base(repo),
		ShortPath: extractShortPath(repo, root),
	}
}

func (r *Repository) Match(re *regexp.Regexp) bool {
	return nil == re.FindStringIndex(r.ShortPath)
}

func extractShortPath(path, root string) string {
	if strings.HasPrefix(path, root) {
		return strings.TrimPrefix(path, root+"/")
	}
	return path
}
