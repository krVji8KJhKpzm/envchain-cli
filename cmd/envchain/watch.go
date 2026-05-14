package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"envchain/internal/env"
	"envchain/internal/keychain"
)

var watchInterval time.Duration

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch <project> <VAR> [VAR...]",
		Short: "Watch variables for changes and print events",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runWatch,
	}
	watchCmd.Flags().DurationVarP(&watchInterval, "interval", "i", 5*time.Second, "polling interval")
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	project := args[0]
	vars := args[1:]

	kc, err := keychain.New("envchain")
	if err != nil {
		return fmt.Errorf("keychain: %w", err)
	}
	store := env.New(kc)

	w := env.NewWatcher(store, watchInterval)
	for _, v := range vars {
		if err := w.Add(project, v); err != nil {
			return fmt.Errorf("watch add %s: %w", v, err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Watching %s/%s (interval: %s)\n", project, v, watchInterval)
	}

	ch := w.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-ch:
			if !ok {
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "CHANGED %s/%s: %q -> %q\n",
				ev.Project, ev.VarName, ev.OldValue, ev.NewValue)
		case <-sig:
			w.Stop()
			return nil
		}
	}
}
