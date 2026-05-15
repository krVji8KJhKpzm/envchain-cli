package env

import (
	"testing"
)

func newQuotaStore(t *testing.T, project string) *QuotaStore {
	t.Helper()
	kc := newTestKeychain()
	qs, err := NewQuotaStore(kc, project)
	if err != nil {
		t.Fatalf("NewQuotaStore: %v", err)
	}
	return qs
}

func TestQuotaDefaultLimit(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	limit, err := qs.GetLimit()
	if err != nil {
		t.Fatalf("GetLimit: %v", err)
	}
	if limit != defaultQuotaLimit {
		t.Errorf("expected default %d, got %d", defaultQuotaLimit, limit)
	}
}

func TestQuotaSetAndGet(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	if err := qs.SetLimit(50); err != nil {
		t.Fatalf("SetLimit: %v", err)
	}
	limit, err := qs.GetLimit()
	if err != nil {
		t.Fatalf("GetLimit: %v", err)
	}
	if limit != 50 {
		t.Errorf("expected 50, got %d", limit)
	}
}

func TestQuotaSetInvalidLimit(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	if err := qs.SetLimit(0); err == nil {
		t.Error("expected error for zero limit")
	}
	if err := qs.SetLimit(-5); err == nil {
		t.Error("expected error for negative limit")
	}
}

func TestCheckQuotaWithinLimit(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	_ = qs.SetLimit(10)
	if err := qs.CheckQuota(5, 3); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckQuotaExceeded(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	_ = qs.SetLimit(10)
	if err := qs.CheckQuota(9, 2); err == nil {
		t.Error("expected quota exceeded error")
	}
}

func TestQuotaRemoveLimit(t *testing.T) {
	qs := newQuotaStore(t, "myproject")
	_ = qs.SetLimit(20)
	if err := qs.RemoveLimit(); err != nil {
		t.Fatalf("RemoveLimit: %v", err)
	}
	limit, err := qs.GetLimit()
	if err != nil {
		t.Fatalf("GetLimit after remove: %v", err)
	}
	if limit != defaultQuotaLimit {
		t.Errorf("expected default %d after remove, got %d", defaultQuotaLimit, limit)
	}
}

func TestNewQuotaStoreEmptyProject(t *testing.T) {
	kc := newTestKeychain()
	_, err := NewQuotaStore(kc, "")
	if err == nil {
		t.Error("expected error for empty project")
	}
}
