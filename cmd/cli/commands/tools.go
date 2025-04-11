// Package commands provides CLI commands for the agent runtime
package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

// NewToolsCommand creates a new command for managing Sam Sepiol tools
func NewToolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage Sam Sepiol tools",
		Long:  `Manage Sam Sepiol tools, including listing, installing, and running tools.`,
	}

	cmd.AddCommand(newToolsListCommand())
	cmd.AddCommand(newToolsRunCommand())
	cmd.AddCommand(newToolsInstallCommand())

	return cmd
}

func newToolsListCommand() *cobra.Command {
	var format string
	var configPath string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available tools",
		Long:  `List all available tools for Sam Sepiol.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			toolsConfig, err := loadToolsConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load tools configuration: %w", err)
			}

			toolsList := toolsConfig.Tools

			switch format {
			case "json":
				jsonData, err := json.MarshalIndent(toolsList, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal tools to JSON: %w", err)
				}
				fmt.Println(string(jsonData))
			case "table":
				fmt.Println("NAME\tDESCRIPTION")
				for _, tool := range toolsList {
					name := tool["function"].(map[string]interface{})["name"].(string)
					description := tool["function"].(map[string]interface{})["description"].(string)
					fmt.Printf("%s\t%s\n", name, description)
				}
			default:
				return fmt.Errorf("unsupported format: %s", format)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (json, table)")
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to tools configuration file")

	return cmd
}

func newToolsRunCommand() *cobra.Command {
	var configPath string
	var paramsStr string

	cmd := &cobra.Command{
		Use:   "run [tool]",
		Short: "Run a tool",
		Long:  `Run a specific tool with the provided parameters.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]

			toolsConfig, err := loadToolsConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load tools configuration: %w", err)
			}

			var toolDef map[string]interface{}
			for _, tool := range toolsConfig.Tools {
				name := tool["function"].(map[string]interface{})["name"].(string)
				if name == toolName {
					toolDef = tool
					break
				}
			}

			if toolDef == nil {
				return fmt.Errorf("tool not found: %s", toolName)
			}

			var params map[string]interface{}
			if paramsStr != "" {
				if err := json.Unmarshal([]byte(paramsStr), &params); err != nil {
					return fmt.Errorf("failed to parse parameters: %w", err)
				}
			} else {
				params = make(map[string]interface{})
			}

			fmt.Printf("Running tool %s with parameters: %v\n", toolName, params)

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to tools configuration file")
	cmd.Flags().StringVarP(&paramsStr, "params", "p", "", "Tool parameters as JSON")

	return cmd
}

func newToolsInstallCommand() *cobra.Command {
	var configPath string
	var force bool

	cmd := &cobra.Command{
		Use:   "install [tool]",
		Short: "Install a tool",
		Long:  `Install a specific tool or all tools if no tool is specified.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolsConfig, err := loadToolsConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load tools configuration: %w", err)
			}

			if len(args) > 0 {
				toolName := args[0]
				fmt.Printf("Installing tool: %s\n", toolName)
			} else {
				fmt.Println("Installing all tools...")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to tools configuration file")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force reinstallation of tools")

	return cmd
}

func loadToolsConfig(configPath string) (*tools.ToolConfig, error) {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		configPath = filepath.Join(homeDir, ".sam", "tools.json")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = filepath.Join(homeDir, ".config", "sam", "tools.json")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				configPath = "tools.json"
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return nil, fmt.Errorf("tools configuration file not found")
				}
			}
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tools configuration file: %w", err)
	}

	var toolsConfig tools.ToolConfig
	if err := json.Unmarshal(data, &toolsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse tools configuration: %w", err)
	}

	return &toolsConfig, nil
}
