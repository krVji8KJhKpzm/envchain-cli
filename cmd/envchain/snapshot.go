package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage project variable snapshots",
	}

	takeCmd := &cobra.Command{
		Use:   "take <project>",
		Short: "Take a snapshot of all variables in a project",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshotTake,
	}
	takeCmd.Flags().String("label", "", "Optional label for the snapshot")

	listCmd := &cobra.Command{
		Use:   "list <project>",
		Short: "List snapshots for a project",
		Args:  cobra.ExactArgs(1),
		RunE:  runSnapshotList,
	}

	restoreCmd := &cobra.Command{
		Use:   "restore <project> <index>",
		Short: "Restore a project to a previous snapshot",
		Args:  cobra.ExactArgs(2),
		RunE:  runSnapshotRestore,
	}

	snapshotCmd.AddCommand(takeCmd, listCmd, restoreCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func newSnapshotStore() (*env.SnapshotStore, error) {
	kc, err := keychain.New("envchain")
	if err != nil {
		return nil, err
	}
	s := env.New(kc)
	return env.NewSnapshotStore(s), nil
}

func runSnapshotTake(cmd *cobra.Command, args []string) error {
	project := args[0]
	label, _ := cmd.Flags().GetString("label")
	ss, err := newSnapshotStore()
	if err != nil {
		return err
	}
	snap, err := ss.Take(project, label)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot taken at %s (%d vars)\n",
		snap.CreatedAt.Format("2006-01-02 15:04:05"), len(snap.Vars))
	return nil
}

func runSnapshotList(cmd *cobra.Command, args []string) error {
	project := args[0]
	ss, err := newSnapshotStore()
	if err != nil {
		return err
	}
	snaps := ss.List(project)
	if len(snaps) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No snapshots found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "INDEX\tCREATED\tLABEL\tVARS")
	for i, s := range snaps {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\n", i,
			s.CreatedAt.Format("2006-01-02 15:04:05"), s.Label, len(s.Vars))
	}
	return w.Flush()
}

func runSnapshotRestore(cmd *cobra.Command, args []string) error {
	project := args[0]
	idx, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("invalid index %q: %w", args[1], err)
	}
	ss, err := newSnapshotStore()
	if err != nil {
		return err
	}
	if err := ss.Restore(project, idx); err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), "Snapshot restored.")
	return nil
}
