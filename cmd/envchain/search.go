package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search variables by name or value across projects",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearch,
	}
	searchCmd.Flags().StringSliceP("projects", "p", nil, "Projects to search (required)")
	searchCmd.Flags().Bool("by-value", false, "Search by value instead of name")
	_ = searchCmd.MarkFlagRequired("projects")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	projects, err := cmd.Flags().GetStringSlice("projects")
	if err != nil {
		return err
	}
	byValue, err := cmd.Flags().GetBool("by-value")
	if err != nil {
		return err
	}

	kc, err := keychain.New("envchain")
	if err != nil {
		return fmt.Errorf("keychain init: %w", err)
	}
	store := env.New(kc)
	searcher := env.NewSearcher(store)

	var results []env.SearchResult
	if byValue {
		results, err = searcher.SearchByValue(projects, query)
	} else {
		results, err = searcher.SearchByName(projects, query)
	}
	if err != nil {
		return fmt.Errorf("search: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stderr, "no matches found")
		return nil
	}

	for _, r := range results {
		fmt.Printf("%-20s %-30s %s\n",
			r.Project,
			r.Name,
			maskValue(r.Value),
		)
	}
	return nil
}

// maskValue partially obscures a secret value for display.
func maskValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}
