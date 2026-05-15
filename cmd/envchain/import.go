package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

var importSkipExisting bool

func init() {
	importCmd := &cobra.Command{
		Use:   "import <project> <file>",
		Short: "Import variables from a .env file into the keychain",
		Long: `Reads KEY=VALUE pairs from a .env file and stores them
under the given project in the OS keychain.

Lines beginning with # and blank lines are ignored.
Quoted values have surrounding double-quotes stripped.`,
		Args:    cobra.ExactArgs(2),
		RunE:    runImport,
		Example: "  envchain import myapp .env\n  envchain import --skip-existing myapp staging.env",
	}
	importCmd.Flags().BoolVar(&importSkipExisting, "skip-existing", false,
		"do not overwrite variables that already exist in the keychain")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	project := args[0]
	filePath := args[1]

	// Verify the file exists and is readable before opening the keychain,
	// so we fail fast with a clear message if the path is wrong.
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("import: file not found: %s", filePath)
		}
		return fmt.Errorf("import: cannot access file %s: %w", filePath, err)
	}

	kc, err := keychain.New()
	if err != nil {
		return fmt.Errorf("keychain: %w", err)
	}
	s := env.New(kc)

	res, err := env.ImportFromFile(s, project, filePath, importSkipExisting)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	for _, e := range res.Errors {
		fmt.Fprintf(os.Stderr, "warning: %s\n", e)
	}
	for _, k := range res.Skipped {
		fmt.Fprintf(cmd.OutOrStdout(), "skipped  %s\n", k)
	}
	for _, k := range res.Imported {
		fmt.Fprintf(cmd.OutOrStdout(), "imported %s\n", k)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nDone: %d imported, %d skipped, %d errors\n",
		len(res.Imported), len(res.Skipped), len(res.Errors))
	return nil
}
