package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envchain/envchain-cli/internal/env"
	"github.com/envchain/envchain-cli/internal/keychain"
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template <project> <template-string>",
		Short: "Render a template string with project env vars",
		Long: `Render a template string substituting {{VAR_NAME}} placeholders
with values stored in the keychain for the given project.

Example:
  envchain template myapp "host={{DB_HOST}} port={{DB_PORT}}"
  envchain template myapp --file config.tmpl`,
		Args: cobra.RangeArgs(1, 2),
		RunE: runTemplate,
	}

	templateCmd.Flags().StringP("file", "f", "", "Read template from file instead of argument")
	templateCmd.Flags().BoolP("list", "l", false, "List placeholders found in the template")
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	project := args[0]

	filePath, _ := cmd.Flags().GetString("file")
	listOnly, _ := cmd.Flags().GetBool("list")

	src, err := resolveTemplateSource(args, filePath)
	if err != nil {
		return err
	}

	if listOnly {
		names := env.ListPlaceholders(src)
		if len(names) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "(no placeholders found)")
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
		return nil
	}

	kc, err := keychain.New(project)
	if err != nil {
		return fmt.Errorf("opening keychain: %w", err)
	}
	store := env.New(kc)

	renderer, err := env.NewTemplateRenderer(store, project)
	if err != nil {
		return err
	}

	out, err := renderer.Render(src)
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), out)
	return nil
}

// resolveTemplateSource returns the template string from either the --file flag
// or the inline argument. It returns an error if neither is provided.
// When a file path is given, it takes precedence over any inline argument.
func resolveTemplateSource(args []string, filePath string) (string, error) {
	if filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("reading template file %q: %w", filePath, err)
		}
		return string(data), nil
	}
	if len(args) < 2 {
		return "", fmt.Errorf("provide a template string as the second argument or use --file")
	}
	if strings.TrimSpace(args[1]) == "" {
		return "", fmt.Errorf("template string must not be empty")
	}
	return args[1], nil
}
