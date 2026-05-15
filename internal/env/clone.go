package env

import "fmt"

// Cloner copies all variables from one project to another.
type Cloner struct {
	store Store
}

// NewCloner creates a Cloner backed by the given Store.
func NewCloner(s Store) *Cloner {
	return &Cloner{store: s}
}

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	Copied  []string
	Skipped []string
}

// CloneProject copies all variables from src project into dst project.
// If overwrite is false, variables that already exist in dst are skipped.
func (c *Cloner) CloneProject(src, dst string, overwrite bool) (*CloneResult, error) {
	if src == "" {
		return nil, fmt.Errorf("source project must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("destination project must not be empty")
	}
	if src == dst {
		return nil, fmt.Errorf("source and destination projects must differ")
	}

	vars, err := c.store.List(src)
	if err != nil {
		return nil, fmt.Errorf("list %q: %w", src, err)
	}

	result := &CloneResult{}

	for _, name := range vars {
		val, err := c.store.Get(src, name)
		if err != nil {
			return nil, fmt.Errorf("get %q/%q: %w", src, name, err)
		}

		if !overwrite {
			existing, err := c.store.Get(dst, name)
			if err == nil && existing != "" {
				result.Skipped = append(result.Skipped, name)
				continue
			}
		}

		if err := c.store.Set(dst, name, val); err != nil {
			return nil, fmt.Errorf("set %q/%q: %w", dst, name, err)
		}
		result.Copied = append(result.Copied, name)
	}

	return result, nil
}
