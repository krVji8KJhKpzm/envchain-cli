package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
)

var ttlStore = env.NewTTLStore()

func init() {
	ttlSetCmd.Flags().StringP("duration", "d", "24h", "TTL duration (e.g. 1h, 30m, 7d)")

	ttlCmd.AddCommand(ttlSetCmd)
	ttlCmd.AddCommand(ttlListCmd)
	tlCmd.AddCommand(ttlCheckCmd)
	rootCmd.AddCommand(ttlCmd)
}

var ttlCmd = &cobra.Command{
	Use:   "ttl",
	Short: "Manage TTL (time-to-live) for environment variables",
}

var ttlSetCmd = &cobra.Command{
	Use:   "set <project> <var>",
	Short: "Set a TTL on an environment variable",
	Args:  cobra.ExactArgs(2),
	RunE:  runTTLSet,
}

var ttlListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all variables with TTLs, highlighting expired ones",
	Args:  cobra.NoArgs,
	RunE:  runTTLList,
}

var ttlCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Exit with non-zero status if any variables are expired",
	Args:  cobra.NoArgs,
	RunE:  runTTLCheck,
}

func runTTLSet(cmd *cobra.Command, args []string) error {
	project, varName := args[0], args[1]
	durStr, _ := cmd.Flags().GetString("duration")

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", durStr, err)
	}

	ttlStore.Set(project, varName, dur)
	fmt.Fprintf(cmd.OutOrStdout(), "TTL set: %s/%s expires in %s\n", project, varName, dur)
	return nil
}

func runTTLList(cmd *cobra.Command, _ []string) error {
	expired := ttlStore.Expired()
	if len(expired) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No TTL entries registered.")
		return nil
	}
	for _, e := range expired {
		fmt.Fprintln(cmd.OutOrStdout(), e.String())
	}
	return nil
}

func runTTLCheck(_ *cobra.Command, _ []string) error {
	expired := ttlStore.Expired()
	if len(expired) == 0 {
		return nil
	}
	for _, e := range expired {
		fmt.Fprintf(os.Stderr, "EXPIRED: %s\n", e.String())
	}
	os.Exit(1)
	return nil
}
