package env

import "fmt"

// Renamer provides variable renaming functionality within a project.
type Renamer struct {
	store Store
}

// NewRenamer creates a Renamer backed by the given Store.
func NewRenamer(store Store) (*Renamer, error) {
	if store == nil {
		return nil, fmt.Errorf("store must not be nil")
	}
	return &Renamer{store: store}, nil
}

// RenameVar renames oldName to newName within project.
// It returns an error if oldName does not exist, newName already exists,
// or either name is invalid.
func (r *Renamer) RenameVar(project, oldName, newName string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if err := validateVarName(oldName); err != nil {
		return fmt.Errorf("invalid source variable name: %w", err)
	}
	if err := validateVarName(newName); err != nil {
		return fmt.Errorf("invalid destination variable name: %w", err)
	}
	if oldName == newName {
		return fmt.Errorf("source and destination names are identical")
	}

	// Ensure source exists.
	value, err := r.store.Get(project, oldName)
	if err != nil {
		return fmt.Errorf("source variable %q not found in project %q: %w", oldName, project, err)
	}

	// Ensure destination does not exist.
	if _, err := r.store.Get(project, newName); err == nil {
		return fmt.Errorf("destination variable %q already exists in project %q", newName, project)
	}

	// Write new name then remove old name.
	if err := r.store.Set(project, newName, value); err != nil {
		return fmt.Errorf("failed to set %q: %w", newName, err)
	}
	if err := r.store.Delete(project, oldName); err != nil {
		// Best-effort rollback.
		_ = r.store.Delete(project, newName)
		return fmt.Errorf("failed to delete %q after rename: %w", oldName, err)
	}
	return nil
}
