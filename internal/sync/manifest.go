// Package sync provides team synchronization support for envchain projects.
// It manages a manifest file that tracks variable names (not values) so teams
// can share which variables are required for a project without exposing secrets.
package sync

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const ManifestFileName = ".envchain.json"

// Manifest represents the shareable project manifest stored in version control.
// It contains variable names and metadata, never secret values.
type Manifest struct {
	Version   int                 `json:"version"`
	Project   string              `json:"project"`
	UpdatedAt time.Time           `json:"updated_at"`
	Vars      map[string]VarMeta  `json:"vars"`
}

// VarMeta holds metadata about an environment variable.
type VarMeta struct {
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
}

// LoadManifest reads and parses a manifest from the given directory.
func LoadManifest(dir string) (*Manifest, error) {
	path := filepath.Join(dir, ManifestFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrManifestNotFound
		}
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Save writes the manifest to the given directory as JSON.
func (m *Manifest) Save(dir string) error {
	m.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, ManifestFileName)
	return os.WriteFile(path, append(data, '\n'), 0644)
}

// NewManifest creates an empty manifest for the given project.
func NewManifest(project string) *Manifest {
	return &Manifest{
		Version: 1,
		Project: project,
		Vars:    make(map[string]VarMeta),
	}
}

// ErrManifestNotFound is returned when no manifest file exists in the directory.
var ErrManifestNotFound = errors.New("sync: manifest not found; run 'envchain init' first")
