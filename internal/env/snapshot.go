package env

import (
	"fmt"
	"sort"
	"time"
)

// Snapshot represents a point-in-time capture of all variables in a project.
type Snapshot struct {
	Project   string            `json:"project"`
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label,omitempty"`
	Vars      map[string]string `json:"vars"`
}

// SnapshotStore manages snapshots for projects.
type SnapshotStore struct {
	store  Store
	snaps  map[string][]Snapshot // keyed by project
}

// NewSnapshotStore creates a SnapshotStore backed by the given Store.
func NewSnapshotStore(s Store) *SnapshotStore {
	return &SnapshotStore{
		store: s,
		snaps: make(map[string][]Snapshot),
	}
}

// Take captures the current state of all variables in project and stores it
// under an optional label.
func (ss *SnapshotStore) Take(project, label string) (Snapshot, error) {
	if project == "" {
		return Snapshot{}, fmt.Errorf("project name must not be empty")
	}
	infos, err := ListVars(ss.store, project)
	if err != nil {
		return Snapshot{}, fmt.Errorf("list vars: %w", err)
	}
	vars := make(map[string]string, len(infos))
	for _, info := range infos {
		val, err := ss.store.Get(project, info.Name)
		if err != nil {
			return Snapshot{}, fmt.Errorf("get %s: %w", info.Name, err)
		}
		vars[info.Name] = val
	}
	snap := Snapshot{
		Project:   project,
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Vars:      vars,
	}
	ss.snaps[project] = append(ss.snaps[project], snap)
	return snap, nil
}

// List returns all snapshots for a project, oldest first.
func (ss *SnapshotStore) List(project string) []Snapshot {
	snaps := ss.snaps[project]
	out := make([]Snapshot, len(snaps))
	copy(out, snaps)
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out
}

// Restore overwrites the project's variables with those from the snapshot at
// index idx (0-based, oldest first).
func (ss *SnapshotStore) Restore(project string, idx int) error {
	snaps := ss.List(project)
	if idx < 0 || idx >= len(snaps) {
		return fmt.Errorf("snapshot index %d out of range (have %d)", idx, len(snaps))
	}
	snap := snaps[idx]
	for k, v := range snap.Vars {
		if err := ss.store.Set(project, k, v); err != nil {
			return fmt.Errorf("restore %s: %w", k, err)
		}
	}
	return nil
}

// Delete removes the snapshot at index idx (0-based, oldest first) for the
// given project. Returns an error if the index is out of range.
func (ss *SnapshotStore) Delete(project string, idx int) error {
	snaps := ss.List(project)
	if idx < 0 || idx >= len(snaps) {
		return fmt.Errorf("snapshot index %d out of range (have %d)", idx, len(snaps))
	}
	// List returns a sorted copy; rebuild the stored slice without the target.
	result := make([]Snapshot, 0, len(snaps)-1)
	for i, s := range snaps {
		if i != idx {
			result = append(result, s)
		}
	}
	ss.snaps[project] = result
	return nil
}
