package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kled",
	Short: "Kled.io Framework CLI",
	Long:  `Kled.io Framework CLI is a command line interface for the Kled.io Framework.`,
}

func init() {
	// Add version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of Kled.io Framework",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Kled.io Framework v0.1.0")
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Add service command
	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage Kled.io Framework services",
		Long:  `Manage Kled.io Framework services including starting, stopping, and listing services.`,
	}
	rootCmd.AddCommand(serviceCmd)

	// Add cluster command
	clusterCmd := &cobra.Command{
		Use:   "cluster",
		Short: "Manage Kled.io Framework clusters",
		Long:  `Manage Kled.io Framework clusters including creating, deleting, and listing clusters.`,
	}
	rootCmd.AddCommand(clusterCmd)

	// Add space command
	spaceCmd := &cobra.Command{
		Use:   "space",
		Short: "Manage Kled.io Framework spaces",
		Long:  `Manage Kled.io Framework spaces including creating, deleting, and listing spaces.`,
	}
	rootCmd.AddCommand(spaceCmd)

	// Add policy command
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage Kled.io Framework policies",
		Long:  `Manage Kled.io Framework policies including creating, deleting, and listing policies.`,
	}
	rootCmd.AddCommand(policyCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
