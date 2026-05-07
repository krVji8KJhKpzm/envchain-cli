package env

import (
	"fmt"
	"regexp"
	"sort"
)

var validTagRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// TagStore manages tags associated with environment variables within a project.
type TagStore struct {
	store Store
}

// NewTagStore creates a TagStore backed by the given Store.
func NewTagStore(s Store) (*TagStore, error) {
	if s == nil {
		return nil, fmt.Errorf("store must not be nil")
	}
	return &TagStore{store: s}, nil
}

func tagKey(varName, tag string) string {
	return fmt.Sprintf("__tag__%s__%s", varName, tag)
}

// AddTag associates a tag with a variable in the given project.
func (t *TagStore) AddTag(project, varName, tag string) error {
	if !validTagRe.MatchString(tag) {
		return fmt.Errorf("invalid tag %q: only alphanumeric, underscore and hyphen allowed", tag)
	}
	return t.store.Set(project, tagKey(varName, tag), "1")
}

// RemoveTag removes a tag from a variable in the given project.
func (t *TagStore) RemoveTag(project, varName, tag string) error {
	return t.store.Delete(project, tagKey(varName, tag))
}

// ListTags returns all tags for the given variable in the project.
func (t *TagStore) ListTags(project, varName string) ([]string, error) {
	prefix := fmt.Sprintf("__tag__%s__", varName)
	vars, err := t.store.List(project)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, v := range vars {
		if len(v.Name) > len(prefix) && v.Name[:len(prefix)] == prefix {
			tags = append(tags, v.Name[len(prefix):])
		}
	}
	sort.Strings(tags)
	return tags, nil
}

// FindByTag returns all variable names that carry the given tag in the project.
func (t *TagStore) FindByTag(project, tag string) ([]string, error) {
	vars, err := t.store.List(project)
	if err != nil {
		return nil, err
	}
	suffix := fmt.Sprintf("__%s", tag)
	prefix := "__tag__"
	var names []string
	for _, v := range vars {
		if len(v.Name) > len(prefix) && v.Name[:len(prefix)] == prefix {
			if len(v.Name) >= len(suffix) && v.Name[len(v.Name)-len(suffix):] == suffix {
				middle := v.Name[len(prefix) : len(v.Name)-len(suffix)]
				names = append(names, middle)
			}
		}
	}
	sort.Strings(names)
	return names, nil
}
