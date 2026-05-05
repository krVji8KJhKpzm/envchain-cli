package env

import (
	"fmt"
	"strings"

	"github.com/envchain-cli/internal/keychain"
)

const servicePrefix = "envchain"

// Store manages environment variables for a named project.
type Store struct {
	kc      *keychain.Keychain
	project string
}

// New creates a new Store for the given project name.
func New(project string, kc *keychain.Keychain) (*Store, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &Store{kc: kc, project: project}, nil
}

// serviceKey returns the keychain service name for this project.
func (s *Store) serviceKey() string {
	return fmt.Sprintf("%s.%s", servicePrefix, s.project)
}

// Set stores an environment variable value in the keychain.
func (s *Store) Set(name, value string) error {
	if err := validateName(name); err != nil {
		return err
	}
	return s.kc.Set(s.serviceKey(), name, value)
}

// Get retrieves an environment variable value from the keychain.
func (s *Store) Get(name string) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}
	return s.kc.Get(s.serviceKey(), name)
}

// Delete removes an environment variable from the keychain.
func (s *Store) Delete(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	return s.kc.Delete(s.serviceKey(), name)
}

// validateName ensures an environment variable name is non-empty and valid.
func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	if strings.ContainsAny(name, " \t\n") {
		return fmt.Errorf("variable name must not contain whitespace")
	}
	return nil
}
