// Package sync implements team synchronization support for envchain projects.
//
// The sync package manages a human-readable manifest file (.envchain.json)
// that is safe to commit to version control. The manifest records:
//
//   - The project name (matching the keychain namespace)
//   - Variable names and their metadata (description, required flag)
//   - A schema version and last-updated timestamp
//
// Crucially, the manifest never stores secret values — only the names and
// metadata needed for teammates to know which variables they must populate
// in their own local keychain.
//
// Typical workflow:
//
//  1. Developer runs 'envchain set PROJECT VAR value' to store a secret.
//  2. Developer runs 'envchain sync' to update .envchain.json with the var name.
//  3. .envchain.json is committed and pushed.
//  4. Teammate pulls and runs 'envchain sync --check' to see missing vars.
//  5. Teammate runs 'envchain set PROJECT VAR <their-value>' to populate them.
package sync
