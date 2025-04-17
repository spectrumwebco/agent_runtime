package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spectrumwebco/agent_runtime/internal/containers"
	"github.com/spectrumwebco/agent_runtime/internal/terminal"
	"github.com/spectrumwebco/agent_runtime/pkg/librechat"
)

func NewNeovimCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neovim",
		Short: "Neovim terminal management",
		Long:  `Commands for managing Neovim terminals in the agent runtime environment.`,
	}

	cmd.AddCommand(newNeovimStartCommand())
	cmd.AddCommand(newNeovimStopCommand())
	cmd.AddCommand(newNeovimListCommand())
	cmd.AddCommand(newNeovimExecCommand())
	cmd.AddCommand(newNeovimBulkCommand())
	cmd.AddCommand(newNeovimCodeCommand())
	cmd.AddCommand(newNeovimProviderCommand())

	return cmd
}

func newNeovimStartCommand() *cobra.Command {
	var useKata bool
	var apiURL string
	var kataConfig string
	var kataRuntime string
	var cpus int
	var memory int
	var debug bool
	var environment string

	cmd := &cobra.Command{
		Use:   "start [id]",
		Short: "Start a Neovim terminal",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			
			options := map[string]interface{}{
				"api_url": apiURL,
				"cpus":    cpus,
				"memory":  memory,
				"debug":   debug,
			}
			
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			if libreChatAPIKey != "" {
				options["librechat_api_key"] = libreChatAPIKey
			}
			
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			terminalType := "neovim"
			if useKata {
				terminalType = "neovim-kata"
				options["kata_config_path"] = kataConfig
				options["kata_runtime_dir"] = kataRuntime
				fmt.Printf("Starting Neovim terminal %s in Kata container\n", id)
			} else {
				fmt.Printf("Starting Neovim terminal %s\n", id)
			}
			
			switch environment {
			case "windows":
				options["os"] = "windows"
				options["shell"] = "powershell.exe"
			case "mac":
				options["os"] = "darwin"
				options["shell"] = "/bin/zsh"
			case "ovhcloud":
				options["provider"] = "ovhcloud"
				options["region"] = os.Getenv("OVHCLOUD_REGION")
			case "flyio":
				options["provider"] = "flyio"
				options["region"] = os.Getenv("FLYIO_REGION")
			case "gcp":
				options["provider"] = "gcp"
				options["project"] = os.Getenv("GCP_PROJECT")
				options["zone"] = os.Getenv("GCP_ZONE")
			case "azure":
				options["provider"] = "azure"
				options["subscription"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
				options["resource_group"] = os.Getenv("AZURE_RESOURCE_GROUP")
			case "aws":
				options["provider"] = "aws"
				options["region"] = os.Getenv("AWS_REGION")
			case "remote":
				options["provider"] = "remote"
				options["host"] = os.Getenv("SSH_HOST")
				options["user"] = os.Getenv("SSH_USERNAME")
				options["key_path"] = os.Getenv("SSH_KEY")
			}
			
			term, err := terminalManager.CreateTerminal(context.Background(), terminalType, id, options)
			if err != nil {
				fmt.Printf("Error creating terminal: %v\n", err)
				os.Exit(1)
			}
			
			if err := term.Start(context.Background()); err != nil {
				fmt.Printf("Error starting terminal: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Terminal %s started successfully\n", id)
		},
	}

	cmd.Flags().BoolVar(&useKata, "kata", false, "Use Kata container")
	cmd.Flags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API URL for Neovim server")
	cmd.Flags().StringVar(&kataConfig, "kata-config", "/etc/kata-containers/configuration.toml", "Kata configuration file")
	cmd.Flags().StringVar(&kataRuntime, "kata-runtime", "/var/run/kata-containers", "Kata runtime directory")
	cmd.Flags().IntVar(&cpus, "cpus", 2, "Number of CPUs to allocate to the container")
	cmd.Flags().IntVar(&memory, "memory", 2048, "Amount of memory in MB to allocate to the container")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode")
	cmd.Flags().StringVar(&environment, "env", "linux", "Environment to run in (linux, windows, mac, ovhcloud, flyio, gcp, azure, aws, remote)")

	return cmd
}

func newNeovimStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop [id]",
		Short: "Stop a Neovim terminal",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			
			apiURL := "http://localhost:8080"
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			term, err := terminalManager.GetTerminal(id)
			if err != nil {
				fmt.Printf("Error getting terminal: %v\n", err)
				os.Exit(1)
			}
			
			if err := term.Stop(context.Background()); err != nil {
				fmt.Printf("Error stopping terminal: %v\n", err)
				os.Exit(1)
			}
			
			if err := terminalManager.RemoveTerminal(context.Background(), id); err != nil {
				fmt.Printf("Error removing terminal: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Terminal %s stopped successfully\n", id)
		},
	}

	return cmd
}

func newNeovimListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Neovim terminals",
		Run: func(cmd *cobra.Command, args []string) {
			apiURL := "http://localhost:8080"
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			terminals := terminalManager.ListTerminals()
			
			if len(terminals) == 0 {
				fmt.Println("No terminals found")
				return
			}
			
			fmt.Println("ID\t\tTYPE\t\tRUNNING")
			fmt.Println("--\t\t----\t\t-------")
			for _, term := range terminals {
				fmt.Printf("%s\t\t%s\t\t%v\n", term.ID(), term.GetType(), term.IsRunning())
			}
		},
	}

	return cmd
}

func newNeovimExecCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec [id] [command]",
		Short: "Execute a command in a Neovim terminal",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			command := strings.Join(args[1:], " ")
			
			apiURL := "http://localhost:8080"
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			term, err := terminalManager.GetTerminal(id)
			if err != nil {
				fmt.Printf("Error getting terminal: %v\n", err)
				os.Exit(1)
			}
			
			output, err := term.Execute(context.Background(), command)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(output)
		},
	}

	return cmd
}

func newNeovimBulkCommand() *cobra.Command {
	var count int
	var useKata bool
	var apiURL string
	var kataConfig string
	var kataRuntime string
	var cpus int
	var memory int
	var debug bool
	var environment string

	cmd := &cobra.Command{
		Use:   "bulk [count]",
		Short: "Create multiple Neovim terminals",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			countStr := args[0]
			var err error
			count, err = strconv.Atoi(countStr)
			if err != nil {
				fmt.Printf("Invalid count: %v\n", err)
				os.Exit(1)
			}
			
			options := map[string]interface{}{
				"api_url": apiURL,
				"cpus":    cpus,
				"memory":  memory,
				"debug":   debug,
			}
			
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			if libreChatAPIKey != "" {
				options["librechat_api_key"] = libreChatAPIKey
			}
			
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			terminalType := "neovim"
			if useKata {
				terminalType = "neovim-kata"
				options["kata_config_path"] = kataConfig
				options["kata_runtime_dir"] = kataRuntime
				fmt.Printf("Creating %d Neovim terminals in Kata containers\n", count)
			} else {
				fmt.Printf("Creating %d Neovim terminals\n", count)
			}
			
			switch environment {
			case "windows":
				options["os"] = "windows"
				options["shell"] = "powershell.exe"
			case "mac":
				options["os"] = "darwin"
				options["shell"] = "/bin/zsh"
			case "ovhcloud":
				options["provider"] = "ovhcloud"
				options["region"] = os.Getenv("OVHCLOUD_REGION")
			case "flyio":
				options["provider"] = "flyio"
				options["region"] = os.Getenv("FLYIO_REGION")
			case "gcp":
				options["provider"] = "gcp"
				options["project"] = os.Getenv("GCP_PROJECT")
				options["zone"] = os.Getenv("GCP_ZONE")
			case "azure":
				options["provider"] = "azure"
				options["subscription"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
				options["resource_group"] = os.Getenv("AZURE_RESOURCE_GROUP")
			case "aws":
				options["provider"] = "aws"
				options["region"] = os.Getenv("AWS_REGION")
			case "remote":
				options["provider"] = "remote"
				options["host"] = os.Getenv("SSH_HOST")
				options["user"] = os.Getenv("SSH_USERNAME")
				options["key_path"] = os.Getenv("SSH_KEY")
			}
			
			terminals, err := terminalManager.CreateBulkTerminals(context.Background(), terminalType, count, options)
			if err != nil {
				fmt.Printf("Error creating terminals: %v\n", err)
				os.Exit(1)
			}
			
			for _, term := range terminals {
				if err := term.Start(context.Background()); err != nil {
					fmt.Printf("Error starting terminal %s: %v\n", term.ID(), err)
					continue
				}
				fmt.Printf("Terminal %s started successfully\n", term.ID())
				time.Sleep(100 * time.Millisecond)
			}
			
			fmt.Printf("%d terminals created and started successfully\n", len(terminals))
		},
	}

	cmd.Flags().BoolVar(&useKata, "kata", false, "Use Kata container")
	cmd.Flags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API URL for Neovim server")
	cmd.Flags().StringVar(&kataConfig, "kata-config", "/etc/kata-containers/configuration.toml", "Kata configuration file")
	cmd.Flags().StringVar(&kataRuntime, "kata-runtime", "/var/run/kata-containers", "Kata runtime directory")
	cmd.Flags().IntVar(&cpus, "cpus", 2, "Number of CPUs to allocate to the container")
	cmd.Flags().IntVar(&memory, "memory", 2048, "Amount of memory in MB to allocate to the container")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug mode")
	cmd.Flags().StringVar(&environment, "env", "linux", "Environment to run in (linux, windows, mac, ovhcloud, flyio, gcp, azure, aws, remote)")

	return cmd
}

func newNeovimCodeCommand() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "code [id] [code]",
		Short: "Execute code in a Neovim terminal using LibreChat Code Interpreter",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			id := args[0]
			code := strings.Join(args[1:], " ")
			
			apiURL := "http://localhost:8080"
			libreChatURL := os.Getenv("LIBRECHAT_URL")
			if libreChatURL == "" {
				libreChatURL = "https://librechat.ai"
			}
			libreChatAPIKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
			
			terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
			
			output, err := terminalManager.ExecuteCodeInTerminal(context.Background(), id, language, code)
			if err != nil {
				fmt.Printf("Error executing code: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Println(output)
		},
	}

	cmd.Flags().StringVar(&language, "lang", "python", "Programming language to execute code in")

	return cmd
}

func newNeovimProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "List available environment providers for Neovim terminals",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available environment providers:")
			fmt.Println("- linux (default): Local Linux VM")
			fmt.Println("- windows: Windows environment")
			fmt.Println("- mac: macOS environment")
			fmt.Println("- ovhcloud: OVHcloud provider")
			fmt.Println("- flyio: Fly.io provider")
			fmt.Println("- gcp: Google Cloud Platform provider")
			fmt.Println("- azure: Microsoft Azure provider")
			fmt.Println("- aws: Amazon Web Services provider")
			fmt.Println("- remote: Remote server via SSH")
			
			fmt.Println("\nEnvironment variables required for providers:")
			fmt.Println("- ovhcloud: OVHCLOUD_REGION")
			fmt.Println("- flyio: FLYIO_REGION")
			fmt.Println("- gcp: GCP_PROJECT, GCP_ZONE")
			fmt.Println("- azure: AZURE_SUBSCRIPTION_ID, AZURE_RESOURCE_GROUP")
			fmt.Println("- aws: AWS_REGION")
			fmt.Println("- remote: SSH_HOST, SSH_USERNAME, SSH_KEY")
		},
	}

	return cmd
}
