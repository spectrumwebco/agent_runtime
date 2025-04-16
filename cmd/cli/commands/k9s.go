package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	k9sModule "github.com/spectrumwebco/agent_runtime/pkg/modules/k9s"
)

func NewK9sCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k9s",
		Short: "K9s terminal UI for Kubernetes",
		Long:  `K9s provides a terminal UI to interact with your Kubernetes clusters.`,
	}

	cmd.AddCommand(newK9sRunCommand())
	cmd.AddCommand(newK9sInstallCommand())
	cmd.AddCommand(newK9sVersionCommand())
	cmd.AddCommand(newK9sResourceCommand())

	return cmd
}

func newK9sRunCommand() *cobra.Command {
	var (
		kubeconfig string
		namespace  string
		context    string
		readOnly   bool
		headless   bool
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run K9s terminal UI",
		Long:  `Run K9s terminal UI to interact with your Kubernetes clusters.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := k9sModule.NewModule(cfg)
			
			if kubeconfig != "" {
				module.WithKubeconfig(kubeconfig)
			}
			
			if namespace != "" {
				module.WithNamespace(namespace)
			}
			
			if context != "" {
				module.WithContext(context)
			}
			
			module.WithReadOnly(readOnly)
			module.WithHeadless(headless)

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			return module.Run(ctx)
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to use")
	cmd.Flags().StringVarP(&context, "context", "c", "", "Kubernetes context to use")
	cmd.Flags().BoolVar(&readOnly, "readonly", false, "Run in read-only mode")
	cmd.Flags().BoolVar(&headless, "headless", false, "Run in headless mode")

	return cmd
}

func newK9sResourceCommand() *cobra.Command {
	var (
		kubeconfig string
		namespace  string
		context    string
		readOnly   bool
		headless   bool
	)

	cmd := &cobra.Command{
		Use:   "resource [resource]",
		Short: "Run K9s and navigate to a specific resource",
		Long:  `Run K9s terminal UI and navigate directly to a specific resource.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := k9sModule.NewModule(cfg)
			
			if kubeconfig != "" {
				module.WithKubeconfig(kubeconfig)
			}
			
			if namespace != "" {
				module.WithNamespace(namespace)
			}
			
			if context != "" {
				module.WithContext(context)
			}
			
			module.WithReadOnly(readOnly)
			module.WithHeadless(headless)

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			return module.RunWithResource(ctx, args[0])
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to use")
	cmd.Flags().StringVarP(&context, "context", "c", "", "Kubernetes context to use")
	cmd.Flags().BoolVar(&readOnly, "readonly", false, "Run in read-only mode")
	cmd.Flags().BoolVar(&headless, "headless", false, "Run in headless mode")

	return cmd
}

func newK9sInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install K9s",
		Long:  `Install K9s terminal UI for Kubernetes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := k9sModule.NewModule(cfg)
			
			fmt.Println("Installing K9s...")
			if err := module.Install(); err != nil {
				return fmt.Errorf("failed to install K9s: %w", err)
			}
			
			fmt.Println("K9s installed successfully!")
			return nil
		},
	}

	return cmd
}

func newK9sVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show K9s version",
		Long:  `Show the installed K9s version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := k9sModule.NewModule(cfg)
			
			version, err := module.GetVersion()
			if err != nil {
				return fmt.Errorf("failed to get K9s version: %w", err)
			}
			
			fmt.Printf("K9s version: %s\n", version)
			return nil
		},
	}

	return cmd
}
