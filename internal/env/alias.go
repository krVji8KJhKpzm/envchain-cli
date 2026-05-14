package env

import (
	"fmt"
	"regexp"
	"strings"
)

var validAliasRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{0,63}$`)

// AliasStore manages short aliases that map to project/variable pairs.
type AliasStore struct {
	kc keychain
}

const aliasProject = "__envchain_aliases__"

// NewAliasStore creates an AliasStore backed by the provided keychain.
func NewAliasStore(kc keychain) (*AliasStore, error) {
	if kc == nil {
		return nil, fmt.Errorf("alias: keychain must not be nil")
	}
	return &AliasStore{kc: kc}, nil
}

// Set registers an alias name pointing to project/varName.
func (a *AliasStore) Set(alias, project, varName string) error {
	if !validAliasRe.MatchString(alias) {
		return fmt.Errorf("alias: invalid alias name %q", alias)
	}
	if strings.TrimSpace(project) == "" {
		return fmt.Errorf("alias: project must not be empty")
	}
	if strings.TrimSpace(varName) == "" {
		return fmt.Errorf("alias: variable name must not be empty")
	}
	value := project + "/" + varName
	return a.kc.Set(aliasProject, alias, value)
}

// Resolve returns the project and variable name for the given alias.
func (a *AliasStore) Resolve(alias string) (project, varName string, err error) {
	raw, err := a.kc.Get(aliasProject, alias)
	if err != nil {
		return "", "", fmt.Errorf("alias: %q not found", alias)
	}
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("alias: corrupt entry for %q", alias)
	}
	return parts[0], parts[1], nil
}

// Remove deletes the alias.
func (a *AliasStore) Remove(alias string) error {
	if err := a.kc.Delete(aliasProject, alias); err != nil {
		return fmt.Errorf("alias: cannot remove %q: %w", alias, err)
	}
	return nil
}

// List returns all registered aliases.
func (a *AliasStore) List() ([]string, error) {
	keys, err := a.kc.List(aliasProject)
	if err != nil {
		return nil, fmt.Errorf("alias: list failed: %w", err)
	}
	return keys, nil
}
