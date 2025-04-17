package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewOTFCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "otf",
		Short: "Open Terraform Framework commands",
		Long:  `Commands for interacting with the Open Terraform Framework.`,
	}

	cmd.AddCommand(newOTFInitCommand())
	cmd.AddCommand(newOTFApplyCommand())
	cmd.AddCommand(newOTFDestroyCommand())

	return cmd
}

func newOTFInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize OTF in the specified directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Initializing OTF in %s\n", path)
		},
	}
	return cmd
}

func newOTFApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [path]",
		Short: "Apply OTF configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Applying OTF configuration in %s\n", path)
		},
	}
	return cmd
}

func newOTFDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy [path]",
		Short: "Destroy OTF resources",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			fmt.Printf("Destroying OTF resources in %s\n", path)
		},
	}
	return cmd
}
