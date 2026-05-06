package env

import (
	"testing"
)

func TestSearchByName(t *testing.T) {
	store := newMockStore()
	_ = store.Set("proj1", "DATABASE_URL", "postgres://localhost")
	_ = store.Set("proj1", "API_KEY", "secret")
	_ = store.Set("proj2", "DATABASE_HOST", "localhost")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByName([]string{"proj1", "proj2"}, "database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchByNameCaseInsensitive(t *testing.T) {
	store := newMockStore()
	_ = store.Set("proj1", "API_KEY", "abc")
	_ = store.Set("proj1", "api_secret", "xyz")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByName([]string{"proj1"}, "API")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchByNameNoMatch(t *testing.T) {
	store := newMockStore()
	_ = store.Set("proj1", "FOO", "bar")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByName([]string{"proj1"}, "NOTFOUND")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchByValue(t *testing.T) {
	store := newMockStore()
	_ = store.Set("proj1", "DB_URL", "postgres://prod-host/db")
	_ = store.Set("proj1", "CACHE_URL", "redis://prod-host")
	_ = store.Set("proj1", "TOKEN", "abc123")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByValue([]string{"proj1"}, "prod-host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchByValueNoMatch(t *testing.T) {
	store := newMockStore()
	_ = store.Set("proj1", "KEY", "value")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByValue([]string{"proj1"}, "notpresent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchAcrossMultipleProjects(t *testing.T) {
	store := newMockStore()
	_ = store.Set("alpha", "SECRET_KEY", "aaa")
	_ = store.Set("beta", "SECRET_TOKEN", "bbb")
	_ = store.Set("gamma", "OTHER", "ccc")

	searcher := NewSearcher(store)
	results, err := searcher.SearchByName([]string{"alpha", "beta", "gamma"}, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
