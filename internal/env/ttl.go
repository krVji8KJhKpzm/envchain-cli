// Package env provides environment variable management backed by the OS keychain.
package env

import (
	"fmt"
	"time"
)

// TTLEntry holds a variable name along with its expiry metadata.
type TTLEntry struct {
	Project   string
	VarName   string
	ExpiresAt time.Time
}

// IsExpired reports whether the TTL entry has passed its expiry time.
func (e TTLEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// String returns a human-readable representation of the TTL entry.
func (e TTLEntry) String() string {
	if e.IsExpired() {
		return fmt.Sprintf("%s/%s: EXPIRED (was %s)", e.Project, e.VarName, e.ExpiresAt.Format(time.RFC3339))
	}
	remaining := time.Until(e.ExpiresAt).Truncate(time.Second)
	return fmt.Sprintf("%s/%s: expires in %s (at %s)", e.Project, e.VarName, remaining, e.ExpiresAt.Format(time.RFC3339))
}

// TTLStore manages expiry records for environment variables.
type TTLStore struct {
	entries []TTLEntry
}

// NewTTLStore creates an empty TTLStore.
func NewTTLStore() *TTLStore {
	return &TTLStore{}
}

// Set registers or updates the TTL for a variable in a project.
func (s *TTLStore) Set(project, varName string, ttl time.Duration) {
	expiry := time.Now().Add(ttl)
	for i, e := range s.entries {
		if e.Project == project && e.VarName == varName {
			s.entries[i].ExpiresAt = expiry
			return
		}
	}
	s.entries = append(s.entries, TTLEntry{Project: project, VarName: varName, ExpiresAt: expiry})
}

// Get returns the TTLEntry for a variable, and whether it was found.
func (s *TTLStore) Get(project, varName string) (TTLEntry, bool) {
	for _, e := range s.entries {
		if e.Project == project && e.VarName == varName {
			return e, true
		}
	}
	return TTLEntry{}, false
}

// Expired returns all entries that have passed their expiry time.
func (s *TTLStore) Expired() []TTLEntry {
	var out []TTLEntry
	for _, e := range s.entries {
		if e.IsExpired() {
			out = append(out, e)
		}
	}
	return out
}

// Remove deletes the TTL record for a variable.
func (s *TTLStore) Remove(project, varName string) {
	filtered := s.entries[:0]
	for _, e := range s.entries {
		if e.Project != project || e.VarName != varName {
			filtered = append(filtered, e)
		}
	}
	s.entries = filtered
}
