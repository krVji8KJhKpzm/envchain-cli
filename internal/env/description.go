package env

import (
	"errors"
	"fmt"
	"strings"
)

const descriptionPrefix = "__desc__"

// DescriptionStore manages human-readable descriptions for environment variables.
type DescriptionStore struct {
	store Store
}

// NewDescriptionStore creates a new DescriptionStore backed by the given Store.
func NewDescriptionStore(s Store) (*DescriptionStore, error) {
	if s == nil {
		return nil, errors.New("description: store must not be nil")
	}
	return &DescriptionStore{store: s}, nil
}

func descKey(project, varName string) string {
	return fmt.Sprintf("%s%s", descriptionPrefix, varName)
}

// Set stores a description for the given variable in the project.
func (d *DescriptionStore) Set(project, varName, description string) error {
	if strings.TrimSpace(project) == "" {
		return errors.New("description: project must not be empty")
	}
	if strings.TrimSpace(varName) == "" {
		return errors.New("description: variable name must not be empty")
	}
	return d.store.Set(project, descKey(project, varName), description)
}

// Get retrieves the description for the given variable in the project.
// Returns an empty string and no error if no description is set.
func (d *DescriptionStore) Get(project, varName string) (string, error) {
	if strings.TrimSpace(project) == "" {
		return "", errors.New("description: project must not be empty")
	}
	if strings.TrimSpace(varName) == "" {
		return "", errors.New("description: variable name must not be empty")
	}
	val, err := d.store.Get(project, descKey(project, varName))
	if errors.Is(err, ErrNotFound) {
		return "", nil
	}
	return val, err
}

// Remove deletes the description for the given variable in the project.
func (d *DescriptionStore) Remove(project, varName string) error {
	if strings.TrimSpace(project) == "" {
		return errors.New("description: project must not be empty")
	}
	if strings.TrimSpace(varName) == "" {
		return errors.New("description: variable name must not be empty")
	}
	err := d.store.Delete(project, descKey(project, varName))
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	return err
}
