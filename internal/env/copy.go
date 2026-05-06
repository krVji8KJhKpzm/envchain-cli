package env

import "fmt"

// CopyVar copies an environment variable from one project to another.
// If destVar is empty, the same variable name is used in the destination project.
func CopyVar(src Store, dst Store, srcProject, dstProject, srcVar, dstVar string) error {
	if dstVar == "" {
		dstVar = srcVar
	}

	value, err := src.Get(srcProject, srcVar)
	if err != nil {
		return fmt.Errorf("copy: get %s/%s: %w", srcProject, srcVar, err)
	}

	if err := dst.Set(dstProject, dstVar, value); err != nil {
		return fmt.Errorf("copy: set %s/%s: %w", dstProject, dstVar, err)
	}

	return nil
}

// MoveVar copies a variable to a new project/name and deletes the original.
func MoveVar(s Store, srcProject, dstProject, srcVar, dstVar string) error {
	if dstVar == "" {
		dstVar = srcVar
	}

	if err := CopyVar(s, s, srcProject, dstProject, srcVar, dstVar); err != nil {
		return fmt.Errorf("move: %w", err)
	}

	if err := s.Delete(srcProject, srcVar); err != nil {
		return fmt.Errorf("move: delete original %s/%s: %w", srcProject, srcVar, err)
	}

	return nil
}

// Store is the interface required for copy/move operations.
type Store interface {
	Get(project, name string) (string, error)
	Set(project, name, value string) error
	Delete(project, name string) error
}
