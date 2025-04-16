package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/modules/gomicro"
)

func NewGoMicroCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gomicro",
		Short: "Go Micro microservices operations",
		Long:  `Perform operations related to Go Micro microservices.`,
	}

	cmd.AddCommand(newGoMicroCreateCommand())
	cmd.AddCommand(newGoMicroListCommand())
	cmd.AddCommand(newGoMicroRunCommand())

	return cmd
}

func newGoMicroCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new microservice",
		Long:  `Create a new Go Micro microservice with the specified name and options.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("service name is required")
			}

			name := args[0]
			version, _ := cmd.Flags().GetString("version")
			address, _ := cmd.Flags().GetString("address")

			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := gomicro.NewModule(cfg)
			err = module.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize Go Micro module: %w", err)
			}

			metadata := map[string]string{
				"created_at": time.Now().Format(time.RFC3339),
				"created_by": "kled-cli",
			}

			service, err := module.CreateService(name, version, address, metadata)
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}

			fmt.Printf("Service %s created successfully\n", service.Name())
			fmt.Printf("Version: %s\n", service.Version())
			if address != "" {
				fmt.Printf("Address: %s\n", address)
			}

			return nil
		},
	}

	cmd.Flags().String("version", "0.1.0", "Service version")
	cmd.Flags().String("address", "", "Service address")

	return cmd
}

func newGoMicroListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List microservices",
		Long:  `List all registered Go Micro microservices.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := gomicro.NewModule(cfg)
			err = module.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize Go Micro module: %w", err)
			}

			services := module.ListServices()
			if len(services) == 0 {
				fmt.Println("No services registered")
				return nil
			}

			fmt.Println("Registered services:")
			for _, service := range services {
				fmt.Printf("- %s (version: %s)\n", service.Name(), service.Version())
			}

			return nil
		},
	}

	return cmd
}

func newGoMicroRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a microservice",
		Long:  `Run a Go Micro microservice with the specified name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("service name is required")
			}

			name := args[0]
			cfg, err := loadConfig(cmd)
			if err != nil {
				return err
			}

			module := gomicro.NewModule(cfg)
			err = module.Initialize()
			if err != nil {
				return fmt.Errorf("failed to initialize Go Micro module: %w", err)
			}

			service, err := module.GetService(name)
			if err != nil {
				version, _ := cmd.Flags().GetString("version")
				address, _ := cmd.Flags().GetString("address")

				metadata := map[string]string{
					"created_at": time.Now().Format(time.RFC3339),
					"created_by": "kled-cli",
				}

				service, err = module.CreateService(name, version, address, metadata)
				if err != nil {
					return fmt.Errorf("failed to create service: %w", err)
				}
			}

			fmt.Printf("Starting service %s (version: %s)\n", service.Name(), service.Version())

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				sig := <-sigCh
				fmt.Printf("Received signal %s, shutting down...\n", sig)
				cancel()
			}()

			go func() {
				err := service.Run()
				if err != nil {
					fmt.Printf("Service error: %s\n", err)
					cancel()
				}
			}()

			<-ctx.Done()
			fmt.Println("Service stopped")

			return nil
		},
	}

	cmd.Flags().String("version", "0.1.0", "Service version")
	cmd.Flags().String("address", "", "Service address")

	return cmd
}
