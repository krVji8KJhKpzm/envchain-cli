package env

import (
	"fmt"
	"strconv"
)

const (
	defaultPriority = 50
	minPriority     = 1
	maxPriority     = 100
)

// PriorityStore manages numeric priority levels for environment variables.
// Higher values indicate higher priority when resolving conflicts across projects.
type PriorityStore struct {
	kc keychain
}

func priorityKey(project, varName string) string {
	return fmt.Sprintf("priority::%s::%s", project, varName)
}

// NewPriorityStore creates a PriorityStore backed by the given keychain.
func NewPriorityStore(kc keychain) *PriorityStore {
	return &PriorityStore{kc: kc}
}

// Set assigns a priority level (1–100) to a variable in a project.
func (p *PriorityStore) Set(project, varName string, priority int) error {
	if project == "" {
		return fmt.Errorf("project name must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	if priority < minPriority || priority > maxPriority {
		return fmt.Errorf("priority must be between %d and %d, got %d", minPriority, maxPriority, priority)
	}
	return p.kc.Set(priorityKey(project, varName), strconv.Itoa(priority))
}

// Get returns the priority of a variable. Returns defaultPriority if not set.
func (p *PriorityStore) Get(project, varName string) (int, error) {
	val, err := p.kc.Get(priorityKey(project, varName))
	if err != nil {
		if isNotFound(err) {
			return defaultPriority, nil
		}
		return 0, err
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("corrupt priority value for %s/%s: %w", project, varName, err)
	}
	return n, nil
}

// Remove deletes the priority entry for a variable.
func (p *PriorityStore) Remove(project, varName string) error {
	err := p.kc.Delete(priorityKey(project, varName))
	if isNotFound(err) {
		return nil
	}
	return err
}

// Compare returns 1 if a has higher priority than b, -1 if lower, 0 if equal.
// Falls back to defaultPriority on lookup errors.
func (p *PriorityStore) Compare(project, varA, varB string) int {
	pa, _ := p.Get(project, varA)
	pb, _ := p.Get(project, varB)
	switch {
	case pa > pb:
		return 1
	case pa < pb:
		return -1
	default:
		return 0
	}
}
