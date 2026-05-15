package env

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var validCategoryRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// CategoryStore manages variable category assignments within a project.
type CategoryStore struct {
	kc keychain
}

func categoryKey(project, varName string) string {
	return fmt.Sprintf("category::%s::%s", project, varName)
}

// NewCategoryStore creates a new CategoryStore backed by the given keychain.
func NewCategoryStore(kc keychain) *CategoryStore {
	return &CategoryStore{kc: kc}
}

// Set assigns a category to a variable within a project.
func (s *CategoryStore) Set(project, varName, category string) error {
	if project == "" {
		return fmt.Errorf("project must not be empty")
	}
	if varName == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	if !validCategoryRe.MatchString(category) {
		return fmt.Errorf("invalid category %q: must match [a-zA-Z0-9_-]+", category)
	}
	return s.kc.Set(categoryKey(project, varName), category)
}

// Get returns the category assigned to a variable, or empty string if none.
func (s *CategoryStore) Get(project, varName string) (string, error) {
	val, err := s.kc.Get(categoryKey(project, varName))
	if err != nil {
		return "", nil
	}
	return val, nil
}

// Remove deletes the category assignment for a variable.
func (s *CategoryStore) Remove(project, varName string) error {
	return s.kc.Delete(categoryKey(project, varName))
}

// ListByCategory returns all variable names in a project that belong to the given category.
func (s *CategoryStore) ListByCategory(project, category string, vars []string) ([]string, error) {
	var matched []string
	for _, v := range vars {
		cat, err := s.Get(project, v)
		if err != nil {
			continue
		}
		if strings.EqualFold(cat, category) {
			matched = append(matched, v)
		}
	}
	sort.Strings(matched)
	return matched, nil
}
