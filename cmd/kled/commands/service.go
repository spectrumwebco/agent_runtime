package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/service"
	"github.com/spectrumwebco/agent_runtime/internal/microservices/registry"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage Kled.io Framework services",
	Long:  `Manage Kled.io Framework services including starting, stopping, and listing services.`,
}

var startServiceCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Start a Kled.io Framework service",
	Long:  `Start a Kled.io Framework service with the specified name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version, _ := cmd.Flags().GetString("version")
		address, _ := cmd.Flags().GetString("address")
		registryAddrs, _ := cmd.Flags().GetStringSlice("registry")

		reg := registry.NewRegistry(
			registry.WithAddrs(registryAddrs...),
			registry.WithTimeout(time.Second*5),
		)

		srv, err := service.NewService(
			service.WithName(fmt.Sprintf("kled.service.%s", name)),
			service.WithVersion(version),
			service.WithRegistry(reg),
			service.WithAddress(address),
			service.WithMetadata(map[string]string{
				"type": "kled.io",
			}),
		)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}

		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
			<-sigChan
			fmt.Println("Shutting down service...")
			os.Exit(0)
		}()

		fmt.Printf("Starting service %s@%s...\n", name, version)
		return srv.Run()
	},
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kled.io Framework services",
	Long:  `List all Kled.io Framework services registered in the registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		registryAddrs, _ := cmd.Flags().GetStringSlice("registry")

		reg := registry.NewRegistry(
			registry.WithAddrs(registryAddrs...),
			registry.WithTimeout(time.Second*5),
		)

		services, err := reg.ListServices()
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		if len(services) == 0 {
			fmt.Println("No services found")
			return nil
		}

		fmt.Println("Services:")
		for _, srv := range services {
			fmt.Printf("  %s@%s\n", srv.Name, srv.Version)
			for _, node := range srv.Nodes {
				fmt.Printf("    - %s\n", node.Address)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startServiceCmd)
	serviceCmd.AddCommand(listServicesCmd)

	startServiceCmd.Flags().String("version", "latest", "Service version")
	startServiceCmd.Flags().String("address", ":0", "Service address")
	startServiceCmd.Flags().StringSlice("registry", []string{"127.0.0.1:8500"}, "Registry addresses")

	listServicesCmd.Flags().StringSlice("registry", []string{"127.0.0.1:8500"}, "Registry addresses")
}
