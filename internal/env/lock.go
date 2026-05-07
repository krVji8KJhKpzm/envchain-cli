package env

import (
	"fmt"
	"time"
)

// LockEntry represents a locked variable that cannot be modified.
type LockEntry struct {
	Project   string
	VarName   string
	LockedAt  time.Time
	LockedBy  string
}

// LockStore manages variable locks backed by a keychain store.
type LockStore struct {
	store Storer
	project string
}

const lockPrefix = "__lock__"

// NewLockStore creates a new LockStore for the given project.
func NewLockStore(store Storer, project string) (*LockStore, error) {
	if project == "" {
		return nil, fmt.Errorf("project name must not be empty")
	}
	return &LockStore{store: store, project: project}, nil
}

// Lock marks a variable as locked, preventing modification.
func (ls *LockStore) Lock(varName, actor string) error {
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	key := lockPrefix + varName
	entry := fmt.Sprintf("%s|%s", actor, time.Now().UTC().Format(time.RFC3339))
	return ls.store.Set(ls.project, key, entry)
}

// Unlock removes the lock from a variable.
func (ls *LockStore) Unlock(varName string) error {
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	key := lockPrefix + varName
	return ls.store.Delete(ls.project, key)
}

// IsLocked returns true if the variable is currently locked.
func (ls *LockStore) IsLocked(varName string) (bool, error) {
	key := lockPrefix + varName
	_, err := ls.store.Get(ls.project, key)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetLockEntry returns metadata about a locked variable.
func (ls *LockStore) GetLockEntry(varName string) (*LockEntry, error) {
	key := lockPrefix + varName
	val, err := ls.store.Get(ls.project, key)
	if err != nil {
		return nil, err
	}
	var actor, tsStr string
	if _, err := fmt.Sscanf(val, "%s", &actor); err != nil {
		return nil, fmt.Errorf("malformed lock entry")
	}
	// parse pipe-delimited format
	for i, c := range val {
		if c == '|' {
			actor = val[:i]
			tsStr = val[i+1:]
			break
		}
	}
	if tsStr == "" {
		return nil, fmt.Errorf("malformed lock entry: missing timestamp")
	}
	ts, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		return nil, fmt.Errorf("malformed lock entry: invalid timestamp: %w", err)
	}
	return &LockEntry{
		Project:  ls.project,
		VarName:  varName,
		LockedAt: ts,
		LockedBy: actor,
	}, nil
}

// isNotFound is a helper that checks for a not-found sentinel in error messages.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "secret not found" ||
		err.Error() == "variable not found"
}
