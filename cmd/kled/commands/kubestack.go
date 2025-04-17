package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewKubestackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubestack",
		Short: "Kubestack commands",
		Long:  `Commands for interacting with Kubestack for Kubernetes infrastructure.`,
	}

	cmd.AddCommand(newKubestackInitCommand())
	cmd.AddCommand(newKubestackApplyCommand())
	cmd.AddCommand(newKubestackDestroyCommand())

	return cmd
}

func newKubestackInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize Kubestack in the specified directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Initializing Kubestack in %s\n", path)
		},
	}
	return cmd
}

func newKubestackApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply Kubestack configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Applying Kubestack configuration in %s\n", path)
		},
	}
	return cmd
}

func newKubestackDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy [path]",
		Short: "Destroy Kubestack resources",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Destroying Kubestack resources in %s\n", path)
		},
	}
	return cmd
}
