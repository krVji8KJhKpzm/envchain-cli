package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

func newNoteStore() (*env.NoteStore, error) {
	kc, err := keychain.New()
	if err != nil {
		return nil, err
	}
	return env.NewNoteStore(kc)
}

func init() {
	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage notes attached to environment variables",
	}

	setCmd := &cobra.Command{
		Use:   "set <project> <var> <note>",
		Short: "Attach a note to a variable",
		Args:  cobra.ExactArgs(3),
		RunE:  runNoteSet,
	}

	getCmd := &cobra.Command{
		Use:   "get <project> <var>",
		Short: "Retrieve the note for a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runNoteGet,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <project> <var>",
		Short: "Remove the note from a variable",
		Args:  cobra.ExactArgs(2),
		RunE:  runNoteRemove,
	}

	noteCmd.AddCommand(setCmd, getCmd, removeCmd)
	rootCmd.AddCommand(noteCmd)
}

func runNoteSet(cmd *cobra.Command, args []string) error {
	ns, err := newNoteStore()
	if err != nil {
		return err
	}
	if err := ns.Set(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Note set for %s/%s\n", args[0], args[1])
	return nil
}

func runNoteGet(cmd *cobra.Command, args []string) error {
	ns, err := newNoteStore()
	if err != nil {
		return err
	}
	note, err := ns.Get(args[0], args[1])
	if err != nil {
		return err
	}
	if note == "" {
		fmt.Fprintf(os.Stdout, "No note set for %s/%s\n", args[0], args[1])
		return nil
	}
	fmt.Fprintln(os.Stdout, note)
	return nil
}

func runNoteRemove(cmd *cobra.Command, args []string) error {
	ns, err := newNoteStore()
	if err != nil {
		return err
	}
	if err := ns.Remove(args[0], args[1]); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Note removed for %s/%s\n", args[0], args[1])
	return nil
}
