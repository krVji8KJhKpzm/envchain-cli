package sync

import (
	"fmt"
	"time"
)

// Syncer orchestrates push/pull operations between the local keychain
// store and a remote Provider, updating the Manifest on success.
type Syncer struct {
	provider Provider
	manifest *Manifest
}

// NewSyncer creates a Syncer using a named provider and an existing manifest.
func NewSyncer(providerName string, m *Manifest) (*Syncer, error) {
	p, err := Get(providerName)
	if err != nil {
		return nil, err
	}
	return &Syncer{provider: p, manifest: m}, nil
}

// Push uploads vars to the provider and records the sync time in the manifest.
func (s *Syncer) Push(project string, vars map[string]string) error {
	if len(vars) == 0 {
		return fmt.Errorf("push: no variables to sync for project %q", project)
	}
	if err := s.provider.Push(project, vars); err != nil {
		return fmt.Errorf("push: %w", err)
	}
	s.manifest.LastSync = time.Now().UTC()
	s.manifest.Provider = s.provider.Name()
	return nil
}

// Pull downloads vars from the provider and records the sync time in the manifest.
func (s *Syncer) Pull(project string) (map[string]string, error) {
	vars, err := s.provider.Pull(project)
	if err != nil {
		return nil, fmt.Errorf("pull: %w", err)
	}
	s.manifest.LastSync = time.Now().UTC()
	s.manifest.Provider = s.provider.Name()
	return vars, nil
}

// ProviderName returns the name of the underlying provider.
func (s *Syncer) ProviderName() string {
	return s.provider.Name()
}
