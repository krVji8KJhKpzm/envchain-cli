package env

import (
	"fmt"
	"time"
)

// RotateResult holds the outcome of a single variable rotation.
type RotateResult struct {
	Project string
	Var     string
	OldSet  bool
	NewVal  string
	Rotated time.Time
}

// Rotator generates new values for environment variables.
type Rotator interface {
	Rotate(project, name, oldValue string) (string, error)
}

// RotateFunc is a function adapter for Rotator.
type RotateFunc func(project, name, oldValue string) (string, error)

func (f RotateFunc) Rotate(project, name, oldValue string) (string, error) {
	return f(project, name, oldValue)
}

// RotateVar rotates a single variable in the given project using the provided Rotator.
// The old value is retrieved, passed to the rotator, and the new value is stored.
func RotateVar(s Store, project, name string, r Rotator) (RotateResult, error) {
	if err := validateProject(project); err != nil {
		return RotateResult{}, err
	}
	if err := validateVarName(name); err != nil {
		return RotateResult{}, err
	}

	oldVal, err := s.Get(project, name)
	oldSet := true
	if err != nil {
		// If not found, proceed with empty old value
		oldSet = false
		oldVal = ""
	}

	newVal, err := r.Rotate(project, name, oldVal)
	if err != nil {
		return RotateResult{}, fmt.Errorf("rotate %s/%s: %w", project, name, err)
	}

	if err := s.Set(project, name, newVal); err != nil {
		return RotateResult{}, fmt.Errorf("store rotated value for %s/%s: %w", project, name, err)
	}

	return RotateResult{
		Project: project,
		Var:     name,
		OldSet:  oldSet,
		NewVal:  newVal,
		Rotated: time.Now(),
	}, nil
}

// RotateAll rotates all listed variables in the project using the provided Rotator.
func RotateAll(s Store, project string, names []string, r Rotator) ([]RotateResult, error) {
	results := make([]RotateResult, 0, len(names))
	for _, name := range names {
		res, err := RotateVar(s, project, name, r)
		if err != nil {
			return results, err
		}
		results = append(results, res)
	}
	return results, nil
}
