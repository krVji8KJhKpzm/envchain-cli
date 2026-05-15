package env

import (
	"testing"
)

func newDepStore(t *testing.T) (*DependencyStore, Store) {
	t.Helper()
	s := newTestStore(t)
	return NewDependencyStore(s), s
}

func TestDependencySetAndGet(t *testing.T) {
	ds, _ := newDepStore(t)
	err := ds.SetDeps("proj", "DB_URL", []string{"DB_HOST", "DB_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps, err := ds.GetDeps("proj", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 2 || deps[0] != "DB_HOST" || deps[1] != "DB_PORT" {
		t.Fatalf("expected [DB_HOST DB_PORT], got %v", deps)
	}
}

func TestDependencyGetNotFound(t *testing.T) {
	ds, _ := newDepStore(t)
	_, err := ds.GetDeps("proj", "MISSING")
	if err == nil {
		t.Fatal("expected error for missing dependency entry")
	}
}

func TestDependencyRemove(t *testing.T) {
	ds, _ := newDepStore(t)
	_ = ds.SetDeps("proj", "KEY", []string{"A"})
	if err := ds.RemoveDeps("proj", "KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := ds.GetDeps("proj", "KEY")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestDependencySatisfied(t *testing.T) {
	ds, s := newDepStore(t)
	_ = s.Set("proj", "DB_HOST", "localhost")
	_ = s.Set("proj", "DB_PORT", "5432")
	_ = ds.SetDeps("proj", "DB_URL", []string{"DB_HOST", "DB_PORT"})

	ok, missing, err := ds.Satisfied("proj", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || len(missing) != 0 {
		t.Fatalf("expected satisfied, got missing=%v", missing)
	}
}

func TestDependencyUnsatisfied(t *testing.T) {
	ds, s := newDepStore(t)
	_ = s.Set("proj", "DB_HOST", "localhost")
	_ = ds.SetDeps("proj", "DB_URL", []string{"DB_HOST", "DB_PORT"})

	ok, missing, err := ds.Satisfied("proj", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected unsatisfied")
	}
	if len(missing) != 1 || missing[0] != "DB_PORT" {
		t.Fatalf("expected missing=[DB_PORT], got %v", missing)
	}
}

func TestDependencyEmptyProject(t *testing.T) {
	ds, _ := newDepStore(t)
	if err := ds.SetDeps("", "KEY", []string{"A"}); err == nil {
		t.Fatal("expected error for empty project")
	}
}

func TestDependencyEmptyVarName(t *testing.T) {
	ds, _ := newDepStore(t)
	if err := ds.SetDeps("proj", "", []string{"A"}); err == nil {
		t.Fatal("expected error for empty var name")
	}
}
