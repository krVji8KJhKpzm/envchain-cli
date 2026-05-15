package env

import (
	"fmt"
	"strings"
)

// ProtectStore manages write-protection flags for environment variables.
// A protected variable cannot be overwritten or deleted without explicit unprotect.
type ProtectStore struct {
	store Store
	project string
}

const protectPrefix = "__protect__"

// NewProtectStore creates a ProtectStore for the given project.
func NewProtectStore(store Store, project string) (*ProtectStore, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &ProtectStore{store: store, project: project}, nil
}

func protectKey(project, varName string) string {
	return fmt.Sprintf("%s%s/%s", protectPrefix, project, varName)
}

// Protect marks a variable as protected.
func (p *ProtectStore) Protect(varName string) error {
	if strings.TrimSpace(varName) == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return p.store.Set(p.project, protectKey(p.project, varName), "1")
}

// Unprotect removes the protection flag from a variable.
func (p *ProtectStore) Unprotect(varName string) error {
	if strings.TrimSpace(varName) == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return p.store.Delete(p.project, protectKey(p.project, varName))
}

// IsProtected returns true if the variable is currently protected.
func (p *ProtectStore) IsProtected(varName string) (bool, error) {
	val, err := p.store.Get(p.project, protectKey(p.project, varName))
	if err != nil {
		return false, nil
	}
	return val == "1", nil
}

// ListProtected returns the names of all protected variables in the project.
func (p *ProtectStore) ListProtected() ([]string, error) {
	prefix := protectKey(p.project, "")
	vars, err := p.store.List(p.project)
	if err != nil {
		return nil, err
	}
	var protected []string
	for _, v := range vars {
		if strings.HasPrefix(v, prefix) {
			name := strings.TrimPrefix(v, prefix)
			protected = append(protected, name)
		}
	}
	return protected, nil
}
