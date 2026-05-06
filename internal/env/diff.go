package env

import "fmt"

// DiffResult holds the comparison between two sets of environment variables.
type DiffResult struct {
	Added   []string // keys present in new but not in old
	Removed []string // keys present in old but not in new
	Changed []string // keys present in both but with different values
}

// String returns a human-readable summary of the diff.
func (d DiffResult) String() string {
	if len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0 {
		return "no changes"
	}
	out := ""
	for _, k := range d.Added {
		out += fmt.Sprintf("+ %s\n", k)
	}
	for _, k := range d.Removed {
		out += fmt.Sprintf("- %s\n", k)
	}
	for _, k := range d.Changed {
		out += fmt.Sprintf("~ %s\n", k)
	}
	return out
}

// HasChanges returns true if the diff is non-empty.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// DiffProjects compares the environment variables of two projects in the store
// and returns a DiffResult describing the differences.
func DiffProjects(s Store, projectA, projectB string) (DiffResult, error) {
	if projectA == "" || projectB == "" {
		return DiffResult{}, fmt.Errorf("project name must not be empty")
	}

	varsA, err := s.GetAll(projectA)
	if err != nil {
		return DiffResult{}, fmt.Errorf("reading project %q: %w", projectA, err)
	}

	varsB, err := s.GetAll(projectB)
	if err != nil {
		return DiffResult{}, fmt.Errorf("reading project %q: %w", projectB, err)
	}

	var result DiffResult

	for k, vA := range varsA {
		if vB, ok := varsB[k]; !ok {
			result.Removed = append(result.Removed, k)
		} else if vA != vB {
			result.Changed = append(result.Changed, k)
		}
	}

	for k := range varsB {
		if _, ok := varsA[k]; !ok {
			result.Added = append(result.Added, k)
		}
	}

	return result, nil
}
