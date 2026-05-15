package env

import (
	"errors"
	"fmt"
	"strings"
)

// Visibility levels for environment variables.
const (
	VisibilityPublic    = "public"
	VisibilityPrivate   = "private"
	VisibilitySecret    = "secret"
)

var validVisibilities = map[string]bool{
	VisibilityPublic:  true,
	VisibilityPrivate: true,
	VisibilitySecret:  true,
}

type visibilityStore interface {
	Set(project, key, value string) error
	Get(project, key string) (string, error)
	Delete(project, key string) error
	List(project string) (map[string]string, error)
}

// VisibilityStore manages per-variable visibility levels.
type VisibilityStore struct {
	store visibilityStore
}

func visibilityKey(varName string) string {
	return fmt.Sprintf("__visibility__%s", varName)
}

// NewVisibilityStore creates a new VisibilityStore backed by the given store.
func NewVisibilityStore(s visibilityStore) *VisibilityStore {
	return &VisibilityStore{store: s}
}

// Set assigns a visibility level to a variable in the given project.
func (v *VisibilityStore) Set(project, varName, level string) error {
	if project == "" {
		return errors.New("project name must not be empty")
	}
	if varName == "" {
		return errors.New("variable name must not be empty")
	}
	level = strings.ToLower(level)
	if !validVisibilities[level] {
		return fmt.Errorf("invalid visibility level %q: must be one of public, private, secret", level)
	}
	return v.store.Set(project, visibilityKey(varName), level)
}

// Get returns the visibility level of a variable. Defaults to "private" if not set.
func (v *VisibilityStore) Get(project, varName string) (string, error) {
	if project == "" {
		return "", errors.New("project name must not be empty")
	}
	if varName == "" {
		return "", errors.New("variable name must not be empty")
	}
	val, err := v.store.Get(project, visibilityKey(varName))
	if err != nil {
		// Default visibility
		return VisibilityPrivate, nil
	}
	return val, nil
}

// Remove deletes the visibility setting for a variable.
func (v *VisibilityStore) Remove(project, varName string) error {
	if project == "" {
		return errors.New("project name must not be empty")
	}
	return v.store.Delete(project, visibilityKey(varName))
}

// ListByLevel returns all variable names in a project that match the given visibility level.
func (v *VisibilityStore) ListByLevel(project, level string) ([]string, error) {
	if project == "" {
		return nil, errors.New("project name must not be empty")
	}
	level = strings.ToLower(level)
	if !validVisibilities[level] {
		return nil, fmt.Errorf("invalid visibility level %q", level)
	}
	all, err := v.store.List(project)
	if err != nil {
		return nil, err
	}
	prefix := "__visibility__"
	var names []string
	for k, val := range all {
		if strings.HasPrefix(k, prefix) && val == level {
			names = append(names, strings.TrimPrefix(k, prefix))
		}
	}
	return names, nil
}
