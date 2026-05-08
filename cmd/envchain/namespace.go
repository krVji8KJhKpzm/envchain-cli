package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envchain-cli/internal/env"
	"github.com/yourorg/envchain-cli/internal/keychain"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage project namespaces",
}

func init() {
	namespaceCmd.AddCommand(namespaceListCmd)
	namespaceCmd.AddCommand(namespaceParseCmd)
	rootCmd.AddCommand(namespaceCmd)
}

var namespaceListCmd = &cobra.Command{
	Use:   "list <namespace>",
	Short: "List projects within a namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  runNamespaceList,
}

var namespaceParseCmd = &cobra.Command{
	Use:   "parse <namespace/project>",
	Short: "Parse a fully-qualified project key into namespace and project",
	Args:  cobra.ExactArgs(1),
	RunE:  runNamespaceParse,
}

func runNamespaceList(cmd *cobra.Command, args []string) error {
	namespace := args[0]

	kc, err := keychain.New("envchain")
	if err != nil {
		return fmt.Errorf("keychain: %w", err)
	}
	store, err := env.New(kc)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	ns, err := env.NewNamespaceStore(store)
	if err != nil {
		return fmt.Errorf("namespace store: %w", err)
	}

	allProjects := store.Projects()
	projects, err := ns.ListProjects(namespace, allProjects)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		fmt.Fprintf(os.Stderr, "no projects found in namespace %q\n", namespace)
		return nil
	}
	for _, p := range projects {
		fmt.Println(p)
	}
	return nil
}

func runNamespaceParse(cmd *cobra.Command, args []string) error {
	ns, project, err := env.ParseProjectKey(args[0])
	if err != nil {
		return err
	}
	fmt.Printf("namespace: %s\nproject:   %s\n", ns, strings.TrimSpace(project))
	return nil
}
