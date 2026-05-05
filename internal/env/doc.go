// Package env provides a project-scoped environment variable store backed
// by the OS keychain via the internal/keychain package.
//
// Each project has its own isolated namespace within the keychain, identified
// by a service key of the form "envchain.<project>". Variable names are
// validated to be non-empty and free of whitespace before any keychain
// operation is performed.
//
// Example usage:
//
//	kc := keychain.New()
//	store, err := env.New("myproject", kc)
//	if err != nil { ... }
//
//	// Store a secret
//	if err := store.Set("API_KEY", "supersecret"); err != nil { ... }
//
//	// Retrieve it later
//	val, err := store.Get("API_KEY")
//
// The export sub-package helpers (FormatShell, FormatDotenv) allow rendering
// retrieved variables into shell-compatible or .env file formats.
package env
