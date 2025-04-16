package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/kledframework"
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

		// Create the service
		service := kledframework.NewService(kledframework.ServiceConfig{
			Name:    name,
			Version: version,
			Address: address,
		})

		// Add a health check endpoint
		service.RegisterHandler("GET", "/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": name,
				"version": version,
			})
		})

		// Handle graceful shutdown
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
			<-sigChan
			fmt.Println("Shutting down service...")
			service.Stop()
			os.Exit(0)
		}()

		// Start the service
		fmt.Printf("Starting service %s@%s...\n", name, version)
		return service.Start()
	},
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kled.io Framework services",
	Long:  `List all Kled.io Framework services registered in the registry.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		registry := kledframework.NewRegistry()
		services := registry.ListServices()

		if len(services) == 0 {
			fmt.Println("No services found")
			return nil
		}

		fmt.Println("Services:")
		for _, srv := range services {
			fmt.Printf("  %s@%s\n", srv.Name, srv.Version)
			fmt.Printf("    - %s\n", srv.Address)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(serviceCmd)

	serviceCmd.AddCommand(startServiceCmd)
	serviceCmd.AddCommand(listServicesCmd)

	startServiceCmd.Flags().String("version", "latest", "Service version")
	startServiceCmd.Flags().String("address", ":8080", "Service address")
}
