package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewClusterCommand creates a new command for managing virtual kubernetes clusters
func NewClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage virtual kubernetes clusters",
		Long:  `Create and manage virtual kubernetes clusters for production deployment of applications.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new virtual cluster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Creating virtual cluster...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List virtual clusters",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing virtual clusters...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a virtual cluster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deleting virtual cluster...")
		},
	})

	return cmd
}
