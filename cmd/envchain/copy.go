package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

var copyCmd = &cobra.Command{
	Use:   "copy <src-project> <var> <dst-project> [dst-var]",
	Short: "Copy an environment variable to another project",
	Args:  cobra.RangeArgs(3, 4),
	RunE:  runCopy,
}

var moveCmd = &cobra.Command{
	Use:   "move <src-project> <var> <dst-project> [dst-var]",
	Short: "Move an environment variable to another project",
	Args:  cobra.RangeArgs(3, 4),
	RunE:  runMove,
}

func init() {
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(moveCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	srcProject, varName, dstProject := args[0], args[1], args[2]
	dstVar := ""
	if len(args) == 4 {
		dstVar = args[3]
	}

	kc, err := keychain.New(serviceName)
	if err != nil {
		return fmt.Errorf("keychain: %w", err)
	}
	store := env.New(kc)

	if err := env.CopyVar(store, store, srcProject, dstProject, varName, dstVar); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	dest := dstVar
	if dest == "" {
		dest = varName
	}
	fmt.Printf("Copied %s/%s → %s/%s\n", srcProject, varName, dstProject, dest)
	return nil
}

func runMove(cmd *cobra.Command, args []string) error {
	srcProject, varName, dstProject := args[0], args[1], args[2]
	dstVar := ""
	if len(args) == 4 {
		dstVar = args[3]
	}

	kc, err := keychain.New(serviceName)
	if err != nil {
		return fmt.Errorf("keychain: %w", err)
	}
	store := env.New(kc)

	if err := env.MoveVar(store, srcProject, dstProject, varName, dstVar); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	dest := dstVar
	if dest == "" {
		dest = varName
	}
	fmt.Printf("Moved %s/%s → %s/%s\n", srcProject, varName, dstProject, dest)
	return nil
}
