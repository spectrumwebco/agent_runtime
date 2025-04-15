package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSpaceCommand creates a new command for building cloud-native applications
func NewSpaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Build cloud-native applications",
		Long:  `Build cloud-native distributed applications on kubernetes with features like hot reloading and declarative workflows.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing project...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "dev",
		Short: "Start development mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting development mode...")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "deploy",
		Short: "Deploy application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deploying application...")
		},
	})

	return cmd
}
