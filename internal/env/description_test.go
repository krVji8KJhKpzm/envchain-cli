package env

import (
	"testing"
)

func newDescriptionStore(t *testing.T) *DescriptionStore {
	t.Helper()
	s := newTestStore(t)
	ds, err := NewDescriptionStore(s)
	if err != nil {
		t.Fatalf("NewDescriptionStore: %v", err)
	}
	return ds
}

func TestDescriptionSetAndGet(t *testing.T) {
	ds := newDescriptionStore(t)

	if err := ds.Set("myproject", "API_KEY", "The API key for the external service"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := ds.Get("myproject", "API_KEY")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "The API key for the external service" {
		t.Errorf("Get = %q; want %q", got, "The API key for the external service")
	}
}

func TestDescriptionGetNotFound(t *testing.T) {
	ds := newDescriptionStore(t)

	got, err := ds.Get("myproject", "MISSING_VAR")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("Get = %q; want empty string", got)
	}
}

func TestDescriptionRemove(t *testing.T) {
	ds := newDescriptionStore(t)

	_ = ds.Set("myproject", "DB_PASS", "Database password")

	if err := ds.Remove("myproject", "DB_PASS"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	got, err := ds.Get("myproject", "DB_PASS")
	if err != nil {
		t.Fatalf("Get after Remove: %v", err)
	}
	if got != "" {
		t.Errorf("Get after Remove = %q; want empty string", got)
	}
}

func TestDescriptionRemoveNotFound(t *testing.T) {
	ds := newDescriptionStore(t)

	if err := ds.Remove("myproject", "NONEXISTENT"); err != nil {
		t.Errorf("Remove of nonexistent key should not error: %v", err)
	}
}

func TestDescriptionEmptyProject(t *testing.T) {
	ds := newDescriptionStore(t)

	if err := ds.Set("", "API_KEY", "desc"); err == nil {
		t.Error("Set with empty project should return error")
	}
	if _, err := ds.Get("", "API_KEY"); err == nil {
		t.Error("Get with empty project should return error")
	}
	if err := ds.Remove("", "API_KEY"); err == nil {
		t.Error("Remove with empty project should return error")
	}
}

func TestDescriptionEmptyVarName(t *testing.T) {
	ds := newDescriptionStore(t)

	if err := ds.Set("myproject", "", "desc"); err == nil {
		t.Error("Set with empty varName should return error")
	}
	if _, err := ds.Get("myproject", ""); err == nil {
		t.Error("Get with empty varName should return error")
	}
}

func TestDescriptionNilStore(t *testing.T) {
	_, err := NewDescriptionStore(nil)
	if err == nil {
		t.Error("NewDescriptionStore(nil) should return error")
	}
}
