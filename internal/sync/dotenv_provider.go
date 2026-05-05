package sync

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DotenvProviderName = "dotenv"

// DotenvProvider syncs variables via a shared .env file on a local or
// mounted filesystem path (useful for team-shared network drives or
// git-crypt encrypted repos).
type DotenvProvider struct {
	// BaseDir is the directory where per-project .env files are stored.
	BaseDir string
}

func NewDotenvProvider(baseDir string) *DotenvProvider {
	return &DotenvProvider{BaseDir: baseDir}
}

func (d *DotenvProvider) Name() string { return DotenvProviderName }

func (d *DotenvProvider) Push(project string, vars map[string]string) error {
	path := d.filePath(project)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("dotenv push: mkdir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("dotenv push: open: %w", err)
	}
	defer f.Close()
	for k, v := range vars {
		_, err := fmt.Fprintf(f, "%s=%s\n", k, quoteDotenvValue(v))
		if err != nil {
			return fmt.Errorf("dotenv push: write: %w", err)
		}
	}
	return nil
}

func (d *DotenvProvider) Pull(project string) (map[string]string, error) {
	path := d.filePath(project)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("dotenv pull: open: %w", err)
	}
	defer f.Close()
	vars := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		vars[parts[0]] = strings.Trim(parts[1], `"`)
	}
	return vars, scanner.Err()
}

func (d *DotenvProvider) filePath(project string) string {
	return filepath.Join(d.BaseDir, project+".env")
}

func quoteDotenvValue(v string) string {
	if strings.ContainsAny(v, " \t\n") {
		return `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
	}
	return v
}
