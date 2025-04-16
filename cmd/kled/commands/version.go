package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new command for displaying the version
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of Kled.io Framework",
		Long:  `Print the version of Kled.io Framework.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Kled.io Framework v0.1.0")
		},
	}
}
