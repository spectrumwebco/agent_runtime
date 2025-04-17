package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "kled",
	Short: "Kled.io Framework CLI",
	Long:  `Kled.io Framework CLI is a command line interface for the Kled.io Framework.`,
}

func init() {
	RootCmd.AddCommand(NewOTFCommand())
	RootCmd.AddCommand(NewKubestackCommand())
	RootCmd.AddCommand(NewLynxCommand())
	RootCmd.AddCommand(NewTerraformOperatorCommand())
	RootCmd.AddCommand(NewNeovimCommand())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
