package env

import (
	"fmt"
	"strings"
)

// LabelStore manages arbitrary key-value labels attached to environment variables.
type LabelStore struct {
	store Store
}

const labelPrefix = "__label__"

// NewLabelStore creates a new LabelStore backed by the given Store.
func NewLabelStore(s Store) *LabelStore {
	return &LabelStore{store: s}
}

func labelKey(project, varName, labelKey string) string {
	return fmt.Sprintf("%s%s__%s__%s", labelPrefix, project, varName, labelKey)
}

// Set attaches a label key-value pair to a variable.
func (l *LabelStore) Set(project, varName, key, value string) error {
	if project == "" || varName == "" || key == "" {
		return fmt.Errorf("project, varName, and key must not be empty")
	}
	if strings.ContainsAny(key, " \t\n") {
		return fmt.Errorf("label key must not contain whitespace")
	}
	return l.store.Set(project, labelKey(project, varName, key), value)
}

// Get retrieves a label value for a variable.
func (l *LabelStore) Get(project, varName, key string) (string, error) {
	if project == "" || varName == "" || key == "" {
		return "", fmt.Errorf("project, varName, and key must not be empty")
	}
	return l.store.Get(project, labelKey(project, varName, key))
}

// Remove deletes a label from a variable.
func (l *LabelStore) Remove(project, varName, key string) error {
	if project == "" || varName == "" || key == "" {
		return fmt.Errorf("project, varName, and key must not be empty")
	}
	return l.store.Delete(project, labelKey(project, varName, key))
}

// List returns all label key-value pairs for a given variable.
func (l *LabelStore) List(project, varName string) (map[string]string, error) {
	if project == "" || varName == "" {
		return nil, fmt.Errorf("project and varName must not be empty")
	}
	vars, err := l.store.List(project)
	if err != nil {
		return nil, err
	}
	prefix := labelKey(project, varName, "")
	result := make(map[string]string)
	for _, v := range vars {
		if strings.HasPrefix(v.Name, prefix) {
			k := strings.TrimPrefix(v.Name, prefix)
			val, err := l.store.Get(project, v.Name)
			if err != nil {
				continue
			}
			result[k] = val
		}
	}
	return result, nil
}
