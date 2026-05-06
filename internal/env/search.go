package env

import (
	"strings"
)

// SearchResult holds a matched variable with its project context.
type SearchResult struct {
	Project string
	Name    string
	Value   string
}

// Searcher provides variable search functionality across projects.
type Searcher struct {
	store Store
}

// NewSearcher creates a new Searcher backed by the given Store.
func NewSearcher(s Store) *Searcher {
	return &Searcher{store: s}
}

// SearchByName returns all variables whose name contains the given substring
// (case-insensitive) across all provided projects.
func (s *Searcher) SearchByName(projects []string, query string) ([]SearchResult, error) {
	q := strings.ToLower(query)
	var results []SearchResult
	for _, proj := range projects {
		vars, err := s.store.List(proj)
		if err != nil {
			return nil, err
		}
		for _, name := range vars {
			if strings.Contains(strings.ToLower(name), q) {
				val, err := s.store.Get(proj, name)
				if err != nil {
					return nil, err
				}
				results = append(results, SearchResult{Project: proj, Name: name, Value: val})
			}
		}
	}
	return results, nil
}

// SearchByValue returns all variables whose value contains the given substring
// (case-sensitive) across all provided projects.
func (s *Searcher) SearchByValue(projects []string, query string) ([]SearchResult, error) {
	var results []SearchResult
	for _, proj := range projects {
		vars, err := s.store.List(proj)
		if err != nil {
			return nil, err
		}
		for _, name := range vars {
			val, err := s.store.Get(proj, name)
			if err != nil {
				return nil, err
			}
			if strings.Contains(val, query) {
				results = append(results, SearchResult{Project: proj, Name: name, Value: val})
			}
		}
	}
	return results, nil
}
