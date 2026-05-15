package env

import (
	"fmt"
	"strings"
)

const bookmarkPrefix = "bookmark:"

// BookmarkStore manages named bookmarks that point to a project+variable pair.
type BookmarkStore struct {
	kc keychain
}

// NewBookmarkStore returns a BookmarkStore backed by the given keychain.
func NewBookmarkStore(kc keychain) *BookmarkStore {
	return &BookmarkStore{kc: kc}
}

func bookmarkKey(name string) string {
	return bookmarkPrefix + strings.ToLower(name)
}

// Set creates or updates a bookmark with the given name pointing to project/variable.
func (b *BookmarkStore) Set(name, project, variable string) error {
	if name == "" {
		return fmt.Errorf("bookmark name must not be empty")
	}
	if project == "" || variable == "" {
		return fmt.Errorf("project and variable must not be empty")
	}
	val := project + "/" + variable
	return b.kc.Set(bookmarkKey(name), val)
}

// Resolve returns the project and variable a bookmark points to.
func (b *BookmarkStore) Resolve(name string) (project, variable string, err error) {
	if name == "" {
		return "", "", fmt.Errorf("bookmark name must not be empty")
	}
	val, err := b.kc.Get(bookmarkKey(name))
	if err != nil {
		return "", "", fmt.Errorf("bookmark %q not found", name)
	}
	parts := strings.SplitN(val, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("bookmark %q has invalid format", name)
	}
	return parts[0], parts[1], nil
}

// Remove deletes a bookmark by name.
func (b *BookmarkStore) Remove(name string) error {
	if name == "" {
		return fmt.Errorf("bookmark name must not be empty")
	}
	return b.kc.Delete(bookmarkKey(name))
}

// List returns all bookmark names and their targets.
func (b *BookmarkStore) List() ([]BookmarkEntry, error) {
	keys, err := b.kc.List()
	if err != nil {
		return nil, err
	}
	var entries []BookmarkEntry
	for _, k := range keys {
		if !strings.HasPrefix(k, bookmarkPrefix) {
			continue
		}
		name := strings.TrimPrefix(k, bookmarkPrefix)
		val, err := b.kc.Get(k)
		if err != nil {
			continue
		}
		parts := strings.SplitN(val, "/", 2)
		if len(parts) != 2 {
			continue
		}
		entries = append(entries, BookmarkEntry{Name: name, Project: parts[0], Variable: parts[1]})
	}
	return entries, nil
}

// BookmarkEntry holds a single bookmark record.
type BookmarkEntry struct {
	Name     string
	Project  string
	Variable string
}

func (e BookmarkEntry) String() string {
	return fmt.Sprintf("%s -> %s/%s", e.Name, e.Project, e.Variable)
}
