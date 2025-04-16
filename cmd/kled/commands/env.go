package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewEnvCommand creates a new command for managing development workspaces
func NewEnvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage development workspaces",
		Long:  `Create and manage development workspaces with features for deploying machines, kubernetes, building workspaces, and deploying workspaces.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new workspace",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating workspace...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing workspaces...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a workspace",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deleting workspace...")
		},
	})

	return cmd
}
