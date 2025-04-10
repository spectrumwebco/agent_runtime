package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
	"gopkg.in/yaml.v3"
)

func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Sam Sepiol configuration",
		Long:  `Manage Sam Sepiol configuration, including viewing and generating configuration files.`,
	}

	cmd.AddCommand(newConfigViewCommand())
	cmd.AddCommand(newConfigGenerateCommand())

	return cmd
}

func newConfigViewCommand() *cobra.Command {
	var configPath string
	var format string

	cmd := &cobra.Command{
		Use:   "view",
		Short: "View configuration",
		Long:  `View the current Sam Sepiol configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			switch format {
			case "json":
				jsonData, err := json.MarshalIndent(cfg, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal configuration to JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			case "yaml":
				yamlData, err := yaml.Marshal(cfg)
				if err != nil {
					return fmt.Errorf("failed to marshal configuration to YAML: %w", err)
				}
				fmt.Println(string(yamlData))
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (json, yaml)")

	return cmd
}

func newConfigGenerateCommand() *cobra.Command {
	var outputPath string
	var format string
	var force bool

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate configuration",
		Long:  `Generate a default Sam Sepiol configuration file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.DefaultConfig()

			if _, err := os.Stat(outputPath); err == nil && !force {
				return fmt.Errorf("output file already exists: %s (use --force to overwrite)", outputPath)
			}

			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			var data []byte
			var err error
			switch format {
			case "json":
				data, err = json.MarshalIndent(cfg, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal configuration to JSON: %w", err)
				}
			case "yaml":
				data, err = yaml.Marshal(cfg)
				if err != nil {
					return fmt.Errorf("failed to marshal configuration to YAML: %w", err)
				}
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			if err := os.WriteFile(outputPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write configuration to file: %w", err)
			}

			fmt.Printf("Configuration generated at: %s\n", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "config.yaml", "Output path for configuration file")
	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (json, yaml)")
	cmd.Flags().BoolVarP(&force, "force", "F", false, "Force overwrite of existing file")

	return cmd
}
