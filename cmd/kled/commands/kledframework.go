package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/kledframework"
)

var kledframeworkCmd = &cobra.Command{
	Use:   "kledframework",
	Short: "Manage Kled.io Framework services",
	Long:  `Manage Kled.io Framework services for the Kled.io Framework.`,
}

var startKledframeworkCmd = &cobra.Command{
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

func init() {
	RootCmd.AddCommand(kledframeworkCmd)
	kledframeworkCmd.AddCommand(startKledframeworkCmd)

	startKledframeworkCmd.Flags().String("version", "latest", "Service version")
	startKledframeworkCmd.Flags().String("address", ":8080", "Service address")
}
