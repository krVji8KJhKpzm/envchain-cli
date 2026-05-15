package env

import (
	"fmt"
	"strings"
)

const depPrefix = "dep:"

// DependencyStore manages variable dependency declarations within a project.
type DependencyStore struct {
	store Store
}

// NewDependencyStore creates a DependencyStore backed by the given Store.
func NewDependencyStore(s Store) *DependencyStore {
	return &DependencyStore{store: s}
}

func depKey(project, varName string) string {
	return depPrefix + project + ":" + varName
}

// SetDeps records that varName in project depends on the given list of variable names.
func (d *DependencyStore) SetDeps(project, varName string, deps []string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	for _, dep := range deps {
		if dep == "" {
			return fmt.Errorf("dependency name must not be empty")
		}
	}
	value := strings.Join(deps, ",")
	return d.store.Set(project, depKey(project, varName), value)
}

// GetDeps returns the list of variable names that varName depends on.
func (d *DependencyStore) GetDeps(project, varName string) ([]string, error) {
	raw, err := d.store.Get(project, depKey(project, varName))
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return []string{}, nil
	}
	return strings.Split(raw, ","), nil
}

// RemoveDeps deletes the dependency declaration for varName in project.
func (d *DependencyStore) RemoveDeps(project, varName string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return d.store.Delete(project, depKey(project, varName))
}

// Satisfied reports whether all declared dependencies for varName exist in the store.
func (d *DependencyStore) Satisfied(project, varName string) (bool, []string, error) {
	deps, err := d.GetDeps(project, varName)
	if err != nil {
		return false, nil, err
	}
	var missing []string
	for _, dep := range deps {
		_, err := d.store.Get(project, dep)
		if err != nil {
			missing = append(missing, dep)
		}
	}
	return len(missing) == 0, missing, nil
}
