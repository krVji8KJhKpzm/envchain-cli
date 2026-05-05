package env

import (
	"fmt"
	"strings"
)

// ExportEntry holds a single environment variable name and value.
type ExportEntry struct {
	Name  string
	Value string
}

// FormatShell formats a slice of ExportEntry values as POSIX shell export
// statements suitable for eval or sourcing.
func FormatShell(entries []ExportEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "export %s=%q\n", e.Name, e.Value)
	}
	return sb.String()
}

// FormatDotenv formats entries in the .env file convention (KEY=VALUE).
func FormatDotenv(entries []ExportEntry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Name, quoteDotenv(e.Value))
	}
	return sb.String()
}

// quoteDotenv wraps values that contain spaces or special characters in
// double quotes, escaping existing double quotes inside.
func quoteDotenv(v string) string {
	if strings.ContainsAny(v, " \t\n#$\"\'\\`") {
		v = strings.ReplaceAll(v, `"`, `\"`)
		return `"` + v + `"`
	}
	return v
}
