package env

import (
	"regexp"
	"strings"
	"testing"
)

func newValidator() *Validator {
	return NewValidator()
}

func TestValidateNoRule(t *testing.T) {
	v := newValidator()
	if err := v.Validate("MY_VAR", "anything"); err != nil {
		t.Fatalf("expected nil for unregistered var, got %v", err)
	}
}

func TestValidateRequired(t *testing.T) {
	v := newValidator()
	_ = v.Register("API_KEY", ValidationRule{Required: true})

	if err := v.Validate("API_KEY", ""); err == nil {
		t.Fatal("expected error for empty required value")
	}
	if err := v.Validate("API_KEY", "secret"); err != nil {
		t.Fatalf("expected nil for non-empty required value, got %v", err)
	}
}

func TestValidateMinLen(t *testing.T) {
	v := newValidator()
	_ = v.Register("TOKEN", ValidationRule{MinLen: 8})

	if err := v.Validate("TOKEN", "short"); err == nil {
		t.Fatal("expected error for value below MinLen")
	}
	if err := v.Validate("TOKEN", "longenough"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMaxLen(t *testing.T) {
	v := newValidator()
	_ = v.Register("CODE", ValidationRule{MaxLen: 4})

	if err := v.Validate("CODE", "toolong"); err == nil {
		t.Fatal("expected error for value exceeding MaxLen")
	}
	if err := v.Validate("CODE", "ok"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatePattern(t *testing.T) {
	v := newValidator()
	_ = v.Register("PORT", ValidationRule{Pattern: regexp.MustCompile(`^\d+$`)})

	if err := v.Validate("PORT", "abc"); err == nil {
		t.Fatal("expected error for non-matching pattern")
	}
	if err := v.Validate("PORT", "8080"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMultipleViolations(t *testing.T) {
	v := newValidator()
	_ = v.Register("SECRET", ValidationRule{Required: true, MinLen: 10})

	err := v.Validate("SECRET", "")
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) < 2 {
		t.Fatalf("expected at least 2 violations, got %d", len(ve.Violations))
	}
	if !strings.Contains(ve.Error(), "SECRET") {
		t.Error("error message should contain var name")
	}
}

func TestRegisterEmptyVarName(t *testing.T) {
	v := newValidator()
	if err := v.Register("", ValidationRule{Required: true}); err == nil {
		t.Fatal("expected error for empty var name")
	}
}

func TestRuleFor(t *testing.T) {
	v := newValidator()
	_ = v.Register("DB_URL", ValidationRule{Required: true, MinLen: 5})

	rule, ok := v.RuleFor("DB_URL")
	if !ok {
		t.Fatal("expected rule to be found")
	}
	if !rule.Required {
		t.Error("expected Required to be true")
	}
	_, ok = v.RuleFor("MISSING")
	if ok {
		t.Error("expected no rule for unregistered var")
	}
}
