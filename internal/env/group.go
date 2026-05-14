// Package env provides environment variable management backed by a keychain.
package env

import (
	"fmt"
	"sort"
	"strings"
)

// GroupStore manages named groups of variables within a project.
type GroupStore struct {
	kv KVStore
}

const groupPrefix = "__group__"

// NewGroupStore returns a GroupStore backed by the given KVStore.
func NewGroupStore(kv KVStore) (*GroupStore, error) {
	if kv == nil {
		return nil, fmt.Errorf("group: kv store must not be nil")
	}
	return &GroupStore{kv: kv}, nil
}

func groupKey(project, group string) string {
	return fmt.Sprintf("%s%s", groupPrefix, group)
}

// AddToGroup adds a variable name to a named group within a project.
func (g *GroupStore) AddToGroup(project, group, varName string) error {
	if project == "" {
		return fmt.Errorf("group: project must not be empty")
	}
	if group == "" {
		return fmt.Errorf("group: group name must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("group: variable name must not be empty")
	}
	key := groupKey(project, group)
	existing, _ := g.kv.Get(project, key)
	members := parseMembers(existing)
	for _, m := range members {
		if m == varName {
			return nil // already present
		}
	}
	members = append(members, varName)
	sort.Strings(members)
	return g.kv.Set(project, key, strings.Join(members, ","))
}

// RemoveFromGroup removes a variable name from a named group.
func (g *GroupStore) RemoveFromGroup(project, group, varName string) error {
	if project == "" {
		return fmt.Errorf("group: project must not be empty")
	}
	key := groupKey(project, group)
	existing, err := g.kv.Get(project, key)
	if err != nil {
		return fmt.Errorf("group %q not found", group)
	}
	members := parseMembers(existing)
	filtered := members[:0]
	for _, m := range members {
		if m != varName {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) == 0 {
		return g.kv.Delete(project, key)
	}
	return g.kv.Set(project, key, strings.Join(filtered, ","))
}

// ListGroup returns all variable names in the given group.
func (g *GroupStore) ListGroup(project, group string) ([]string, error) {
	if project == "" {
		return nil, fmt.Errorf("group: project must not be empty")
	}
	key := groupKey(project, group)
	val, err := g.kv.Get(project, key)
	if err != nil {
		return []string{}, nil
	}
	return parseMembers(val), nil
}

// ListGroups returns all group names defined in the project.
func (g *GroupStore) ListGroups(project string) ([]string, error) {
	if project == "" {
		return nil, fmt.Errorf("group: project must not be empty")
	}
	all, err := g.kv.List(project)
	if err != nil {
		return nil, err
	}
	var groups []string
	for _, v := range all {
		if strings.HasPrefix(v.Name, groupPrefix) {
			groups = append(groups, strings.TrimPrefix(v.Name, groupPrefix))
		}
	}
	sort.Strings(groups)
	return groups, nil
}

func parseMembers(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
