package env

import (
	"fmt"
	"time"
)

// ExpiryStore manages expiration policies for environment variables.
// Unlike TTL (which is relative), expiry uses absolute timestamps.
type ExpiryStore struct {
	kc keychain
	project string
}

// ExpiryEntry holds the expiration metadata for a variable.
type ExpiryEntry struct {
	VarName   string
	ExpiresAt time.Time
}

// IsExpired reports whether the entry has passed its expiration time.
func (e ExpiryEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// String returns a human-readable representation of the expiry entry.
func (e ExpiryEntry) String() string {
	if e.IsExpired() {
		return fmt.Sprintf("%s: expired at %s", e.VarName, e.ExpiresAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("%s: expires at %s (in %s)", e.VarName, e.ExpiresAt.Format(time.RFC3339), time.Until(e.ExpiresAt).Round(time.Second))
}

func expiryKey(project, varName string) string {
	return fmt.Sprintf("expiry::%s::%s", project, varName)
}

// NewExpiryStore creates an ExpiryStore for the given project.
func NewExpiryStore(kc keychain, project string) (*ExpiryStore, error) {
	if project == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &ExpiryStore{kc: kc, project: project}, nil
}

// SetExpiry records an absolute expiration time for a variable.
func (s *ExpiryStore) SetExpiry(varName string, expiresAt time.Time) error {
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	return s.kc.Set(expiryKey(s.project, varName), expiresAt.UTC().Format(time.RFC3339))
}

// GetExpiry retrieves the expiration entry for a variable.
func (s *ExpiryStore) GetExpiry(varName string) (ExpiryEntry, error) {
	raw, err := s.kc.Get(expiryKey(s.project, varName))
	if err != nil {
		return ExpiryEntry{}, fmt.Errorf("no expiry set for %q: %w", varName, err)
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return ExpiryEntry{}, fmt.Errorf("corrupt expiry value for %q: %w", varName, err)
	}
	return ExpiryEntry{VarName: varName, ExpiresAt: t}, nil
}

// RemoveExpiry deletes the expiration policy for a variable.
func (s *ExpiryStore) RemoveExpiry(varName string) error {
	return s.kc.Delete(expiryKey(s.project, varName))
}

// ListExpired returns all variables whose expiry has passed, given a list of known var names.
func (s *ExpiryStore) ListExpired(varNames []string) ([]ExpiryEntry, error) {
	var expired []ExpiryEntry
	for _, name := range varNames {
		entry, err := s.GetExpiry(name)
		if err != nil {
			continue // no expiry set — skip
		}
		if entry.IsExpired() {
			expired = append(expired, entry)
		}
	}
	return expired, nil
}
