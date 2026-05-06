package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

var rotateCmd = &cobra.Command{
	Use:   "rotate <project> <VAR> [VAR...]",
	Short: "Rotate one or more environment variables using a generator",
	Long: `Rotate re-generates values for the specified variables in the project.

By default, a random 32-character alphanumeric token is generated for each variable.
The old value is overwritten in the OS keychain.`,
	Args: cobra.MinimumNArgs(2),
	RunE: runRotate,
}

var rotateLength int

func init() {
	rotateCmd.Flags().IntVarP(&rotateLength, "length", "l", 32, "length of generated token")
	rootCmd.AddCommand(rotateCmd)
}

func runRotate(cmd *cobra.Command, args []string) error {
	project := args[0]
	vars := args[1:]

	kc, err := keychain.New(appName)
	if err != nil {
		return fmt.Errorf("open keychain: %w", err)
	}
	store := env.New(kc)

	rotator := env.RotateFunc(func(project, name, oldValue string) (string, error) {
		return generateToken(rotateLength), nil
	})

	results, err := env.RotateAll(store, project, vars, rotator)
	if err != nil {
		return err
	}

	for _, r := range results {
		status := "created"
		if r.OldSet {
			status = "rotated"
		}
		fmt.Fprintf(os.Stdout, "%s %s/%s\n", status, r.Project, r.Var)
	}
	return nil
}
