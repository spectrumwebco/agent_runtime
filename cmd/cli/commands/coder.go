package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/devenv/coder"
	codermodule "github.com/spectrumwebco/agent_runtime/pkg/modules/coder"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func NewCoderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coder",
		Short: "Self-hosted cloud development environments",
		Long:  `Coder provides self-hosted cloud development environments defined with Terraform.`,
	}

	cmd.AddCommand(newCoderWorkspaceCommand())
	cmd.AddCommand(newCoderTemplateCommand())
	cmd.AddCommand(newCoderKataCommand())

	return cmd
}

func newCoderWorkspaceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage Coder workspaces",
		Long:  `Manage Coder workspaces for cloud development environments.`,
	}

	cmd.AddCommand(newCoderWorkspaceCreateCommand())
	cmd.AddCommand(newCoderWorkspaceListCommand())
	cmd.AddCommand(newCoderWorkspaceGetCommand())
	cmd.AddCommand(newCoderWorkspaceDeleteCommand())
	cmd.AddCommand(newCoderWorkspaceStartCommand())
	cmd.AddCommand(newCoderWorkspaceStopCommand())

	return cmd
}

func newCoderWorkspaceCreateCommand() *cobra.Command {
	var templateID string
	var paramsFile string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new workspace",
		Long:  `Create a new Coder workspace from a template.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			provisioner, err := coder.NewProvisioner(cfg)
			if err != nil {
				return fmt.Errorf("failed to create provisioner: %v", err)
			}

			var params map[string]interface{}
			if paramsFile != "" {
				paramsData, err := os.ReadFile(paramsFile)
				if err != nil {
					return fmt.Errorf("failed to read params file: %v", err)
				}

				if err := json.Unmarshal(paramsData, &params); err != nil {
					return fmt.Errorf("failed to parse params file: %v", err)
				}
			}

			result, err := provisioner.ProvisionWorkspace(context.Background(), name, templateID, params)
			if err != nil {
				return fmt.Errorf("failed to provision workspace: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("Workspace ID: %s\n", result.WorkspaceID)
				fmt.Printf("Status: %s\n", result.Status)
				fmt.Printf("Started At: %s\n", result.StartedAt.Format(time.RFC3339))
				if !result.CompletedAt.IsZero() {
					fmt.Printf("Completed At: %s\n", result.CompletedAt.Format(time.RFC3339))
				}
				if result.Error != "" {
					fmt.Printf("Error: %s\n", result.Error)
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&templateID, "template", "t", "", "Template ID (required)")
	cmd.Flags().StringVarP(&paramsFile, "params", "p", "", "Path to JSON file containing template parameters")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")
	cmd.MarkFlagRequired("template")

	return cmd
}

func newCoderWorkspaceListCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workspaces",
		Long:  `List all Coder workspaces.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			workspaces, err := client.ListWorkspaces(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list workspaces: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(workspaces, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal workspaces: %v", err)
				}
				fmt.Println(string(json))
			case "table":
				fmt.Printf("%-20s %-20s %-15s %-25s\n", "ID", "NAME", "STATUS", "CREATED AT")
				for _, workspace := range workspaces {
					fmt.Printf("%-20s %-20s %-15s %-25s\n",
						workspace.ID, workspace.Name, workspace.Status,
						workspace.CreatedAt.Format(time.RFC3339))
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (json, table)")

	return cmd
}

func newCoderWorkspaceGetCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a workspace",
		Long:  `Get a Coder workspace by ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			workspace, err := client.GetWorkspace(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to get workspace: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(workspace, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal workspace: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("ID: %s\n", workspace.ID)
				fmt.Printf("Name: %s\n", workspace.Name)
				fmt.Printf("Template ID: %s\n", workspace.TemplateID)
				fmt.Printf("Status: %s\n", workspace.Status)
				fmt.Printf("Created At: %s\n", workspace.CreatedAt.Format(time.RFC3339))
				if !workspace.UpdatedAt.IsZero() {
					fmt.Printf("Updated At: %s\n", workspace.UpdatedAt.Format(time.RFC3339))
				}
				if workspace.URL != "" {
					fmt.Printf("URL: %s\n", workspace.URL)
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")

	return cmd
}

func newCoderWorkspaceDeleteCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a workspace",
		Long:  `Delete a Coder workspace by ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			provisioner, err := coder.NewProvisioner(cfg)
			if err != nil {
				return fmt.Errorf("failed to create provisioner: %v", err)
			}

			result, err := provisioner.DestroyWorkspace(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to destroy workspace: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("Workspace ID: %s\n", result.WorkspaceID)
				fmt.Printf("Status: %s\n", result.Status)
				fmt.Printf("Started At: %s\n", result.StartedAt.Format(time.RFC3339))
				if !result.CompletedAt.IsZero() {
					fmt.Printf("Completed At: %s\n", result.CompletedAt.Format(time.RFC3339))
				}
				if result.Error != "" {
					fmt.Printf("Error: %s\n", result.Error)
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")

	return cmd
}

func newCoderWorkspaceStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [id]",
		Short: "Start a workspace",
		Long:  `Start a Coder workspace by ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			err = client.StartWorkspace(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to start workspace: %v", err)
			}

			fmt.Printf("Workspace %s started\n", id)

			return nil
		},
	}

	return cmd
}

func newCoderWorkspaceStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [id]",
		Short: "Stop a workspace",
		Long:  `Stop a Coder workspace by ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			err = client.StopWorkspace(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to stop workspace: %v", err)
			}

			fmt.Printf("Workspace %s stopped\n", id)

			return nil
		},
	}

	return cmd
}

func newCoderTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage Coder templates",
		Long:  `Manage Coder templates for cloud development environments.`,
	}

	cmd.AddCommand(newCoderTemplateListCommand())
	cmd.AddCommand(newCoderTemplateGetCommand())
	cmd.AddCommand(newCoderTemplateCreateCommand())

	return cmd
}

func newCoderTemplateListCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List templates",
		Long:  `List all Coder templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			templates, err := client.ListTemplates(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list templates: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(templates, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal templates: %v", err)
				}
				fmt.Println(string(json))
			case "table":
				fmt.Printf("%-20s %-20s %-40s %-25s\n", "ID", "NAME", "DESCRIPTION", "CREATED AT")
				for _, template := range templates {
					fmt.Printf("%-20s %-20s %-40s %-25s\n",
						template.ID, template.Name, template.Description,
						template.CreatedAt.Format(time.RFC3339))
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format (json, table)")

	return cmd
}

func newCoderTemplateGetCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a template",
		Long:  `Get a Coder template by ID.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			client, err := coder.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			template, err := client.GetTemplate(context.Background(), id)
			if err != nil {
				return fmt.Errorf("failed to get template: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(template, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal template: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("ID: %s\n", template.ID)
				fmt.Printf("Name: %s\n", template.Name)
				fmt.Printf("Description: %s\n", template.Description)
				fmt.Printf("Created At: %s\n", template.CreatedAt.Format(time.RFC3339))
				if !template.UpdatedAt.IsZero() {
					fmt.Printf("Updated At: %s\n", template.UpdatedAt.Format(time.RFC3339))
				}
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")

	return cmd
}

func newCoderTemplateCreateCommand() *cobra.Command {
	var description string
	var variablesFile string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new template",
		Long:  `Create a new Coder template.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			terraformProvider := coder.NewTerraformProvider(cfg)

			var variables map[string]interface{}
			if variablesFile != "" {
				variablesData, err := os.ReadFile(variablesFile)
				if err != nil {
					return fmt.Errorf("failed to read variables file: %v", err)
				}

				if err := json.Unmarshal(variablesData, &variables); err != nil {
					return fmt.Errorf("failed to parse variables file: %v", err)
				}
			}

			template, err := terraformProvider.CreateTemplate(context.Background(), name, description, variables)
			if err != nil {
				return fmt.Errorf("failed to create template: %v", err)
			}

			switch outputFormat {
			case "json":
				json, err := json.MarshalIndent(template, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal template: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("ID: %s\n", template.ID)
				fmt.Printf("Name: %s\n", template.Name)
				fmt.Printf("Description: %s\n", template.Description)
				fmt.Printf("Created At: %s\n", template.CreatedAt.Format(time.RFC3339))
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Template description")
	cmd.Flags().StringVarP(&variablesFile, "variables", "v", "", "Path to JSON file containing template variables")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")

	return cmd
}

func newCoderKataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kata",
		Short: "Manage Kata containers",
		Long:  `Manage Kata containers for Coder workspaces.`,
	}

	cmd.AddCommand(newCoderKataTemplateCommand())

	return cmd
}

func newCoderKataTemplateCommand() *cobra.Command {
	var description string
	var image string
	var cpu int
	var memory int
	var disk int
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "template [name]",
		Short: "Create a Kata container template",
		Long:  `Create a Coder template for Kata containers.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load("")
			if err != nil {
				return fmt.Errorf("failed to load configuration: %v", err)
			}

			terraformProvider := coder.NewTerraformProvider(cfg)

			config := &coder.KataContainerConfig{
				Image:  image,
				CPU:    cpu,
				Memory: memory,
				Disk:   disk,
			}

			templateDir, err := terraformProvider.GenerateKataContainerTemplate(context.Background(), name, description, config)
			if err != nil {
				return fmt.Errorf("failed to generate template: %v", err)
			}

			switch outputFormat {
			case "json":
				result := struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					TemplateDir string `json:"template_dir"`
					Config      *coder.KataContainerConfig `json:"config"`
				}{
					Name:        name,
					Description: description,
					TemplateDir: templateDir,
					Config:      config,
				}

				json, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(json))
			case "summary":
				fmt.Printf("Name: %s\n", name)
				fmt.Printf("Description: %s\n", description)
				fmt.Printf("Template Directory: %s\n", templateDir)
				fmt.Printf("Image: %s\n", config.Image)
				fmt.Printf("CPU: %d\n", config.CPU)
				fmt.Printf("Memory: %d MB\n", config.Memory)
				fmt.Printf("Disk: %d GB\n", config.Disk)
			default:
				return fmt.Errorf("unknown output format: %s", outputFormat)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&description, "description", "d", "", "Template description")
	cmd.Flags().StringVarP(&image, "image", "i", "ubuntu:20.04", "Container image")
	cmd.Flags().IntVarP(&cpu, "cpu", "c", 2, "CPU cores")
	cmd.Flags().IntVarP(&memory, "memory", "m", 4096, "Memory in MB")
	cmd.Flags().IntVarP(&disk, "disk", "s", 10, "Disk size in GB")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "summary", "Output format (json, summary)")

	return cmd
}
