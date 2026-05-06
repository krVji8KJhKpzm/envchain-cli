//go:build integration

package env_test

import (
	"testing"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func TestCopyVarIntegration(t *testing.T) {
	kc, err := keychain.New("envchain-test-copy")
	if err != nil {
		t.Skipf("keychain unavailable: %v", err)
	}

	store := env.New(kc)

	const src = "integ-copy-src"
	const dst = "integ-copy-dst"
	const varName = "COPY_TEST_VAR"
	const value = "integration-value"

	t.Cleanup(func() {
		_ = store.Delete(src, varName)
		_ = store.Delete(dst, varName)
	})

	if err := store.Set(src, varName, value); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if err := env.CopyVar(store, store, src, dst, varName, ""); err != nil {
		t.Fatalf("CopyVar: %v", err)
	}

	got, err := store.Get(dst, varName)
	if err != nil {
		t.Fatalf("Get dst: %v", err)
	}
	if got != value {
		t.Errorf("expected %q, got %q", value, got)
	}

	// original must still exist
	if _, err := store.Get(src, varName); err != nil {
		t.Error("original deleted after copy")
	}
}

func TestMoveVarIntegration(t *testing.T) {
	kc, err := keychain.New("envchain-test-move")
	if err != nil {
		t.Skipf("keychain unavailable: %v", err)
	}

	store := env.New(kc)

	const src = "integ-move-src"
	const dst = "integ-move-dst"
	const varName = "MOVE_TEST_VAR"
	const value = "move-value"

	t.Cleanup(func() {
		_ = store.Delete(src, varName)
		_ = store.Delete(dst, varName)
	})

	if err := store.Set(src, varName, value); err != nil {
		t.Fatalf("Set: %v", err)
	}

	if err := env.MoveVar(store, src, dst, varName, ""); err != nil {
		t.Fatalf("MoveVar: %v", err)
	}

	if _, err := store.Get(src, varName); err == nil {
		t.Error("original should be gone after move")
	}

	got, err := store.Get(dst, varName)
	if err != nil {
		t.Fatalf("Get dst: %v", err)
	}
	if got != value {
		t.Errorf("expected %q, got %q", value, got)
	}
}
