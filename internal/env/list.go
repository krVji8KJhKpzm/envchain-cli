package env

import (
	"sort"
	"strings"
)

// VarInfo holds metadata about a stored environment variable.
type VarInfo struct {
	Project string
	Name    string
}

// String returns a human-readable representation of the variable reference.
func (v VarInfo) String() string {
	return v.Project + "/" + v.Name
}

// ListVars returns a sorted list of variable names stored for the given project.
// It returns an empty slice if the project has no variables.
func (s *Store) ListVars(project string) ([]string, error) {
	if err := validateProject(project); err != nil {
		return nil, err
	}

	keys, err := s.kc.List(serviceKey(project))
	if err != nil {
		return nil, err
	}

	// Filter to only valid variable names (defensive: skip any corrupted keys).
	var names []string
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k != "" {
			names = append(names, k)
		}
	}

	sort.Strings(names)
	return names, nil
}

// ListAll returns VarInfo entries for every variable across all given projects.
// Projects are processed in the order provided.
func (s *Store) ListAll(projects []string) ([]VarInfo, error) {
	var result []VarInfo
	for _, project := range projects {
		names, err := s.ListVars(project)
		if err != nil {
			return nil, err
		}
		for _, name := range names {
			result = append(result, VarInfo{Project: project, Name: name})
		}
	}
	return result, nil
}
