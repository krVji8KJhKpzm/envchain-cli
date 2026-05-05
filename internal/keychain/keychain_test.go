package keychain_test

import (
	"testing"

	"github.com/zalando/go-keyring"

	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	// Use the in-memory mock keyring for all tests so no real OS keychain is touched.
	keyring.MockInit()
}

func TestSetAndGet(t *testing.T) {
	client := keychain.New("test-project")

	if err := client.Set("API_KEY", "supersecret"); err != nil {
		t.Fatalf("Set: unexpected error: %v", err)
	}

	val, err := client.Get("API_KEY")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if val != "supersecret" {
		t.Errorf("Get: expected %q, got %q", "supersecret", val)
	}
}

func TestGetNotFound(t *testing.T) {
	client := keychain.New("test-project")

	_, err := client.Get("NONEXISTENT_KEY")
	if err == nil {
		t.Fatal("Get: expected error for missing key, got nil")
	}
	if err != keychain.ErrNotFound {
		t.Errorf("Get: expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	client := keychain.New("test-project")

	if err := client.Set("TO_DELETE", "value"); err != nil {
		t.Fatalf("Set: unexpected error: %v", err)
	}
	if err := client.Delete("TO_DELETE"); err != nil {
		t.Fatalf("Delete: unexpected error: %v", err)
	}
	_, err := client.Get("TO_DELETE")
	if err != keychain.ErrNotFound {
		t.Errorf("Get after Delete: expected ErrNotFound, got %v", err)
	}
}

func TestDeleteNotFound(t *testing.T) {
	client := keychain.New("test-project")

	err := client.Delete("DOES_NOT_EXIST")
	if err != keychain.ErrNotFound {
		t.Errorf("Delete: expected ErrNotFound, got %v", err)
	}
}

func TestEmptyKeyErrors(t *testing.T) {
	client := keychain.New("test-project")

	if err := client.Set("", "val"); err == nil {
		t.Error("Set: expected error for empty key")
	}
	if _, err := client.Get(""); err == nil {
		t.Error("Get: expected error for empty key")
	}
	if err := client.Delete(""); err == nil {
		t.Error("Delete: expected error for empty key")
	}
}

func TestProjectIsolation(t *testing.T) {
	a := keychain.New("project-a")
	b := keychain.New("project-b")

	if err := a.Set("SHARED_KEY", "value-a"); err != nil {
		t.Fatalf("Set project-a: %v", err)
	}

	_, err := b.Get("SHARED_KEY")
	if err != keychain.ErrNotFound {
		t.Errorf("project-b should not see project-a keys; got err=%v", err)
	}
}
