package env

import (
	"fmt"
	"strings"
)

// PinStore manages pinned (read-only protected) environment variables.
// A pinned variable cannot be overwritten or deleted without first unpinning it.
type PinStore struct {
	kc keychain
	project string
}

const pinPrefix = "__pin__"

// NewPinStore creates a PinStore for the given project.
func NewPinStore(kc keychain, project string) (*PinStore, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &PinStore{kc: kc, project: project}, nil
}

func pinKey(varName string) string {
	return pinPrefix + varName
}

// Pin marks a variable as pinned.
func (p *PinStore) Pin(varName string) error {
	if strings.TrimSpace(varName) == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return p.kc.Set(p.project, pinKey(varName), "1")
}

// Unpin removes the pin from a variable.
func (p *PinStore) Unpin(varName string) error {
	if strings.TrimSpace(varName) == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return p.kc.Delete(p.project, pinKey(varName))
}

// IsPinned reports whether the variable is currently pinned.
func (p *PinStore) IsPinned(varName string) (bool, error) {
	val, err := p.kc.Get(p.project, pinKey(varName))
	if err != nil {
		// Not found means not pinned.
		return false, nil
	}
	return val == "1", nil
}

// ListPinned returns the names of all pinned variables in the project.
func (p *PinStore) ListPinned() ([]string, error) {
	keys, err := p.kc.List(p.project)
	if err != nil {
		return nil, err
	}
	var pinned []string
	for _, k := range keys {
		if strings.HasPrefix(k, pinPrefix) {
			pinned = append(pinned, strings.TrimPrefix(k, pinPrefix))
		}
	}
	return pinned, nil
}
