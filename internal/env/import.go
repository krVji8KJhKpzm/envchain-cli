// Package env provides environment variable management backed by a secure keychain.
package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ImportResult holds the outcome of a bulk import operation.
type ImportResult struct {
	Imported []string
	Skipped  []string
	Errors   []string
}

// ImportFromDotenv reads key=value pairs from a .env-formatted reader and
// stores each variable under the given project. Existing variables are
// overwritten unless skipExisting is true.
//
// Supported line formats:
//   - KEY=value
//   - KEY="quoted value"
//   - # comment lines (ignored)
//   - blank lines (ignored)
func ImportFromDotenv(s Store, project string, r io.Reader, skipExisting bool) (ImportResult, error) {
	if err := validateProject(project); err != nil {
		return ImportResult{}, err
	}

	var result ImportResult
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx <= 0 {
			result.Errors = append(result.Errors, fmt.Sprintf("malformed line: %q", line))
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.Trim(strings.TrimSpace(line[idx+1:]), `"`)

		if err := validateVarName(key); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", key, err))
			continue
		}

		if skipExisting {
			if _, err := s.Get(project, key); err == nil {
				result.Skipped = append(result.Skipped, key)
				continue
			}
		}

		if err := s.Set(project, key, val); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", key, err))
			continue
		}
		result.Imported = append(result.Imported, key)
	}
	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("reading input: %w", err)
	}
	return result, nil
}

// ImportFromFile is a convenience wrapper that opens the given file path and
// calls ImportFromDotenv.
func ImportFromFile(s Store, project, path string, skipExisting bool) (ImportResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	return ImportFromDotenv(s, project, f, skipExisting)
}

// Summary returns a human-readable one-line description of the import result,
// e.g. "imported 3, skipped 1, errors 0".
func (r ImportResult) Summary() string {
	return fmt.Sprintf("imported %d, skipped %d, errors %d",
		len(r.Imported), len(r.Skipped), len(r.Errors))
}
