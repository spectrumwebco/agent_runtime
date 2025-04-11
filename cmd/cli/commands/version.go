// Package commands provides CLI commands for the agent runtime
package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of Sam Sepiol
	Version   = "0.1.0"
	// BuildDate is the date when Sam Sepiol was built
	BuildDate = "unknown"
	// GitCommit is the git commit hash from which Sam Sepiol was built
	GitCommit = "unknown"
	// GoVersion is the Go version used to build Sam Sepiol
	GoVersion = runtime.Version()
	// Platform is the operating system and architecture Sam Sepiol is running on
	Platform  = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// NewVersionCommand creates a new command for displaying version information
func NewVersionCommand() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version information of Sam Sepiol.`,
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Printf("Version:    %s\n", Version)
				fmt.Printf("Build Date: %s\n", BuildDate)
				fmt.Printf("Git Commit: %s\n", GitCommit)
				fmt.Printf("Go Version: %s\n", GoVersion)
				fmt.Printf("Platform:   %s\n", Platform)
			} else {
				fmt.Printf("Sam Sepiol v%s\n", Version)
			}
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print detailed version information")

	return cmd
}
