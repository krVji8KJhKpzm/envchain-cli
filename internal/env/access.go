package env

import (
	"fmt"
	"strings"
	"time"
)

// AccessEntry records a single read access to an environment variable.
type AccessEntry struct {
	Project   string
	VarName   string
	Actor     string
	AccessedAt time.Time
}

// AccessStore tracks read access events per project/variable.
type AccessStore struct {
	kv KVStore
}

const accessPrefix = "__access__"

func accessKey(project, varName string) string {
	return fmt.Sprintf("%s%s__%s", accessPrefix, project, varName)
}

// NewAccessStore creates an AccessStore backed by the given KVStore.
func NewAccessStore(kv KVStore) *AccessStore {
	return &AccessStore{kv: kv}
}

// Record saves an access entry for the given project and variable.
func (a *AccessStore) Record(project, varName, actor string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("varName must not be empty")
	}
	entry := fmt.Sprintf("%s|%s|%s", actor, varName, time.Now().UTC().Format(time.RFC3339))
	key := accessKey(project, varName)
	existing, _ := a.kv.Get(key)
	lines := []string{}
	if existing != "" {
		lines = strings.Split(existing, "\n")
	}
	lines = append([]string{entry}, lines...)
	const maxEntries = 50
	if len(lines) > maxEntries {
		lines = lines[:maxEntries]
	}
	return a.kv.Set(key, strings.Join(lines, "\n"))
}

// List returns access entries for a given project and variable, newest first.
func (a *AccessStore) List(project, varName string) ([]AccessEntry, error) {
	if project == "" {
		return nil, fmt.Errorf("project must not be empty")
	}
	raw, err := a.kv.Get(accessKey(project, varName))
	if err != nil || raw == "" {
		return []AccessEntry{}, nil
	}
	var entries []AccessEntry
	for _, line := range strings.Split(raw, "\n") {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) != 3 {
			continue
		}
		t, err := time.Parse(time.RFC3339, parts[2])
		if err != nil {
			continue
		}
		entries = append(entries, AccessEntry{
			Project:    project,
			VarName:    parts[1],
			Actor:      parts[0],
			AccessedAt: t,
		})
	}
	return entries, nil
}

// Clear removes all access records for a project/variable pair.
func (a *AccessStore) Clear(project, varName string) error {
	return a.kv.Delete(accessKey(project, varName))
}
