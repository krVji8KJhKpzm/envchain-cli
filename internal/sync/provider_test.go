package sync

import (
	"testing"
)

type mockProvider struct{ name string }

func (m *mockProvider) Name() string                              { return m.name }
func (m *mockProvider) Push(_ string, _ map[string]string) error { return nil }
func (m *mockProvider) Pull(_ string) (map[string]string, error) {
	return map[string]string{"FOO": "bar"}, nil
}

func TestRegisterAndGet(t *testing.T) {
	p := &mockProvider{name: "mock"}
	Register(p)

	got, err := Get("mock")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if got.Name() != "mock" {
		t.Errorf("expected name %q, got %q", "mock", got.Name())
	}
}

func TestGetNotFoundProvider(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	nfe, ok := err.(*ErrProviderNotFound)
	if !ok {
		t.Fatalf("expected *ErrProviderNotFound, got %T", err)
	}
	if nfe.Name != "nonexistent" {
		t.Errorf("unexpected name in error: %s", nfe.Name)
	}
}

func TestAvailable(t *testing.T) {
	Register(&mockProvider{name: "alpha"})
	Register(&mockProvider{name: "beta"})

	names := Available()
	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["alpha"] || !found["beta"] {
		t.Errorf("Available missing expected providers, got: %v", names)
	}
}
