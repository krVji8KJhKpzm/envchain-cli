package env

import (
	"fmt"
	"strings"
)

// InheritStore manages variable inheritance between projects.
// A child project can inherit variables from a parent project,
// with child values taking precedence over parent values.
type InheritStore struct {
	store Store
}

const inheritKey = "__inherit__"

// NewInheritStore creates a new InheritStore backed by the given Store.
func NewInheritStore(s Store) *InheritStore {
	return &InheritStore{store: s}
}

// SetParent sets the parent project for the given child project.
func (h *InheritStore) SetParent(child, parent string) error {
	if strings.TrimSpace(child) == "" {
		return fmt.Errorf("child project must not be empty")
	}
	if strings.TrimSpace(parent) == "" {
		return fmt.Errorf("parent project must not be empty")
	}
	if child == parent {
		return fmt.Errorf("project cannot inherit from itself")
	}
	return h.store.Set(child, inheritKey, parent)
}

// GetParent returns the parent project for the given child, or empty string if none.
func (h *InheritStore) GetParent(child string) (string, error) {
	v, err := h.store.Get(child, inheritKey)
	if err != nil {
		return "", nil
	}
	return v, nil
}

// RemoveParent removes the inheritance link from the given child project.
func (h *InheritStore) RemoveParent(child string) error {
	return h.store.Delete(child, inheritKey)
}

// Resolve returns the effective value for varName in project, walking up
// the inheritance chain. Child values take precedence over parent values.
func (h *InheritStore) Resolve(project, varName string) (string, error) {
	visited := map[string]bool{}
	current := project
	for current != "" {
		if visited[current] {
			return "", fmt.Errorf("inheritance cycle detected at project %q", current)
		}
		visited[current] = true
		v, err := h.store.Get(current, varName)
		if err == nil {
			return v, nil
		}
		parent, perr := h.GetParent(current)
		if perr != nil || parent == "" {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("variable %q not found in project %q or its parents", varName, project)
}
