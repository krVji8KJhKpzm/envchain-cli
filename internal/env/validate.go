package env

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidationRule defines a rule applied to an environment variable value.
type ValidationRule struct {
	Name    string
	Pattern *regexp.Regexp
	MinLen  int
	MaxLen  int
	Required bool
}

// ValidationError holds all violations found during validation.
type ValidationError struct {
	VarName    string
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %q: %s", e.VarName, strings.Join(e.Violations, "; "))
}

// Validator validates env var values against registered rules.
type Validator struct {
	rules map[string]ValidationRule
}

// NewValidator returns a new Validator.
func NewValidator() *Validator {
	return &Validator{rules: make(map[string]ValidationRule)}
}

// Register adds or replaces a validation rule for the given variable name.
func (v *Validator) Register(varName string, rule ValidationRule) error {
	if varName == "" {
		return errors.New("variable name must not be empty")
	}
	v.rules[varName] = rule
	return nil
}

// Validate checks value against the rule registered for varName.
// Returns nil if no rule is registered or all checks pass.
func (v *Validator) Validate(varName, value string) error {
	rule, ok := v.rules[varName]
	if !ok {
		return nil
	}

	var violations []string

	if rule.Required && value == "" {
		violations = append(violations, "value is required")
	}

	if rule.MinLen > 0 && len(value) < rule.MinLen {
		violations = append(violations, fmt.Sprintf("value length %d is below minimum %d", len(value), rule.MinLen))
	}

	if rule.MaxLen > 0 && len(value) > rule.MaxLen {
		violations = append(violations, fmt.Sprintf("value length %d exceeds maximum %d", len(value), rule.MaxLen))
	}

	if rule.Pattern != nil && value != "" && !rule.Pattern.MatchString(value) {
		violations = append(violations, fmt.Sprintf("value does not match pattern %q", rule.Pattern.String()))
	}

	if len(violations) > 0 {
		return &ValidationError{VarName: varName, Violations: violations}
	}
	return nil
}

// RuleFor returns the registered rule for varName and whether it exists.
func (v *Validator) RuleFor(varName string) (ValidationRule, bool) {
	r, ok := v.rules[varName]
	return r, ok
}
