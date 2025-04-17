package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewTerraformOperatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "terraform-operator",
		Short: "Terraform Operator commands",
		Long:  `Commands for interacting with Terraform Operator for Kubernetes.`,
	}

	cmd.AddCommand(newTerraformOperatorDeployCommand())
	cmd.AddCommand(newTerraformOperatorListCommand())
	cmd.AddCommand(newTerraformOperatorApplyCommand())

	return cmd
}

func newTerraformOperatorDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy Terraform Operator to the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Deploying Terraform Operator to the cluster")
		},
	}
	return cmd
}

func newTerraformOperatorListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Terraform resources managed by the operator",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing Terraform resources managed by the operator")
		},
	}
	return cmd
}

func newTerraformOperatorApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [name]",
		Short: "Apply a Terraform resource",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("Applying Terraform resource %s\n", name)
		},
	}
	return cmd
}
