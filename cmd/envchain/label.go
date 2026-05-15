package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	labelCmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels on environment variables",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <var> <key> <value>",
		Short: "Attach a label to an environment variable",
		Args:  cobra.ExactArgs(4),
		RunE:  runLabelSet,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <var> <key>",
		Short: "Remove a label from an environment variable",
		Args:  cobra.ExactArgs(3),
		RunE:  runLabelRemove,
	}

	listCmd := &cobra.Command{
		Use:   "list <project> <var>",
		Short: "List all labels on an environment variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runLabelList,
	}

	labelCmd.AddCommand(setCmd, removeCmd, listCmd)
	rootCmd.AddCommand(labelCmd)
}

func newLabelStore() (*env.LabelStore, error) {
	kc, err := keychain.New("envchain")
	if err != nil {
		return nil, fmt.Errorf("keychain init: %w", err)
	}
	s := env.New(kc)
	return env.NewLabelStore(s), nil
}

func runLabelSet(cmd *cobra.Command, args []string) error {
	ls, err := newLabelStore()
	if err != nil {
		return err
	}
	if err := ls.Set(args[0], args[1], args[2], args[3]); err != nil {
		return fmt.Errorf("label set: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Label %q set on %s/%s\n", args[2], args[0], args[1])
	return nil
}

func runLabelRemove(cmd *cobra.Command, args []string) error {
	ls, err := newLabelStore()
	if err != nil {
		return err
	}
	if err := ls.Remove(args[0], args[1], args[2]); err != nil {
		return fmt.Errorf("label remove: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Label %q removed from %s/%s\n", args[2], args[0], args[1])
	return nil
}

func runLabelList(cmd *cobra.Command, args []string) error {
	ls, err := newLabelStore()
	if err != nil {
		return err
	}
	labels, err := ls.List(args[0], args[1])
	if err != nil {
		return fmt.Errorf("label list: %w", err)
	}
	if len(labels) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No labels for %s/%s\n", args[0], args[1])
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "KEY\tVALUE")
	for k, v := range labels {
		fmt.Fprintf(w, "%s\t%s\n", k, v)
	}
	return w.Flush()
}
