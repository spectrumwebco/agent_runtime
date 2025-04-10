package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/server"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func NewServeCommand() *cobra.Command {
	var configPath string
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Sam Sepiol API server",
		Long:  `Start the Sam Sepiol API server, which provides a REST API for interacting with the agent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			if port != 0 {
				cfg.Server.Port = port
			}
			if host != "" {
				cfg.Server.Host = host
			}

			srv, err := server.NewServer(cfg)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			go func() {
				if err := srv.Start(); err != nil {
					fmt.Printf("Failed to start server: %v\n", err)
					os.Exit(1)
				}
			}()

			fmt.Printf("Server started on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			fmt.Println("Shutting down server...")
			if err := srv.Stop(); err != nil {
				return fmt.Errorf("failed to stop server: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().IntVarP(&port, "port", "p", 0, "Port to listen on (overrides config)")
	cmd.Flags().StringVarP(&host, "host", "H", "", "Host to listen on (overrides config)")

	return cmd
}
