# Coder Integration for Kled.io Framework

## Overview

This document details the integration of [Coder](https://github.com/coder/coder) into the Kled.io Framework. Coder is a self-hosted platform that enables organizations to provision remote development environments through infrastructure as code. This integration provides a seamless way for users to create, manage, and use cloud development environments within the Kled.io ecosystem.

## Key Components

The Coder integration consists of the following key components:

1. **Client API**: Provides programmatic access to Coder functionality
2. **Provisioner**: Manages workspace provisioning and lifecycle
3. **Terraform Provider**: Enables infrastructure as code for development environments
4. **Kata Container Support**: Specialized container runtime for enhanced security and isolation
5. **CLI Commands**: User-friendly command-line interface for Coder operations

## Comprehensive File Mapping

The following table provides a detailed mapping between the original Coder repository files and their corresponding implementations in the agent_runtime repository:

| Coder File | agent_runtime File | Description |
|------------|-------------------|-------------|
| agent/agentcontainers/containers.go | pkg/interfaces/container.go | Container interfaces and data structures |
| agent/agentcontainers/containers_dockercli.go | internal/containers/podman.go | Container implementation (Podman instead of Docker) |
| agent/agentcontainers/api.go | internal/mcp/librechat_provider.go | API endpoints for container management |
| agent/agentcontainers/devcontainer.go | internal/containers/kata.go | Container environment configuration |
| pty/terminal.go | pkg/interfaces/terminal.go | Terminal interface definitions |
| agent/agentssh/agentssh.go | internal/terminal/neovim.go | Terminal implementation |
| agent/agent.go | internal/terminal/manager.go | Terminal management |
| cli/agent.go | cmd/kled/commands/neovim.go | CLI commands for terminal management |
| cli/cliui/agent.go | cmd/kled/commands/neovim.go | CLI UI components |
| N/A | pkg/librechat/client.go | LibreChat Code Interpreter API client |
| N/A | internal/mcp/librechat_provider.go | LibreChat MCP provider |
| N/A | internal/containers/kata.go | Kata container integration |
| N/A | internal/terminal/neovim_kata.go | Neovim in Kata container |
| N/A | internal/exec/exec.go | Command execution utilities |

## Component Mapping

The following table maps the components from the Coder repository to their corresponding implementations in the agent_runtime repository:

| Coder Component | agent_runtime Component | Notes |
|-----------------|-------------------------|-------|
| Docker CLI | Podman CLI | Replaced Docker with Podman for container management |
| Reconnecting PTY | Neovim Terminal | Using Neovim as the terminal implementation |
| SSH Server | Terminal Manager | Extended to support multiple terminal types |
| Container Management | Container Interface | Abstracted to support multiple container runtimes |
| N/A | LibreChat Integration | Added new component for code execution |
| N/A | Kata Container Integration | Added new component for sandboxed environments |
| N/A | Multi-Environment Support | Added support for Windows, Mac, and cloud providers |

## Environment Provider Mapping

The following table maps the supported environments in the agent_runtime implementation:

| Environment | agent_runtime Implementation | Configuration |
|-------------|----------------------------|---------------|
| Linux VM | Default environment | Native support |
| Windows | Windows-specific options | Uses PowerShell |
| Mac | Mac-specific options | Uses zsh |
| OVHcloud | OVHcloud provider | Requires OVHCLOUD_REGION |
| Fly.io | Fly.io provider | Requires FLYIO_REGION |
| GCP | GCP provider | Requires GCP_PROJECT, GCP_ZONE |
| Azure | Azure provider | Requires AZURE_SUBSCRIPTION_ID, AZURE_RESOURCE_GROUP |
| AWS | AWS provider | Requires AWS_REGION |
| Remote Server | SSH provider | Requires SSH_HOST, SSH_USERNAME, SSH_KEY |

## Removed Components

The following components from Coder were not ported to the agent_runtime implementation:

- Frontend application (web UI)
- Helm charts for Kubernetes deployment
- coderd server (main API server)
- DERP networking (direct connection protocol)
- Database migrations and models
- User authentication and authorization
- Workspace provisioning
- Template management
- Audit logging
- Metrics collection
- OAuth providers
- SCIM integration
- External authentication providers

## File Count Comparison

| Directory | Coder Files | agent_runtime Files | Notes |
|-----------|------------|---------------------|-------|
| agent/ | 15 | N/A | Replaced with internal/ structure |
| cli/ | 8 | 1 | Simplified CLI implementation |
| pty/ | 3 | N/A | Integrated into terminal package |
| internal/ | N/A | 7 | New structure for agent_runtime |
| pkg/ | N/A | 3 | New structure for agent_runtime |
| cmd/ | N/A | 1 | New structure for agent_runtime |

## Directory Structure

```
agent_runtime/
├── internal/
│   ├── containers/
│   │   ├── podman.go         # Podman container implementation
│   │   └── kata.go           # Kata container integration
│   ├── terminal/
│   │   ├── neovim.go         # Neovim terminal implementation
│   │   ├── neovim_kata.go    # Neovim in Kata container
│   │   └── manager.go        # Terminal management
│   ├── exec/
│   │   └── exec.go           # Command execution utilities
│   └── mcp/
│       └── librechat_provider.go # LibreChat MCP provider
├── pkg/
│   ├── interfaces/
│   │   ├── container.go      # Container interfaces
│   │   └── terminal.go       # Terminal interfaces
│   └── librechat/
│       └── client.go         # LibreChat Code Interpreter API client
└── cmd/
    └── kled/
        └── commands/
            ├── neovim.go     # Neovim CLI commands
            └── root.go       # Root CLI command
```

## Configuration

The Coder integration is configured through the main configuration file. The following settings are available:

```yaml
coder:
  url: "http://localhost:3000"    # URL of the Coder server
  token: "your-token-here"        # Authentication token
  timeout: 30                     # Request timeout in seconds
kata:
  url: "http://localhost:8080"    # URL of the Kata Containers API
  timeout: 30                     # Request timeout in seconds
```

## Usage

### CLI Commands

The Coder integration provides a comprehensive set of CLI commands:

```bash
# List workspaces
kled neovim list

# Create a new terminal
kled neovim start my-terminal

# Execute a command in a terminal
kled neovim exec my-terminal "ls -la"

# Stop a terminal
kled neovim stop my-terminal

# Create multiple terminals
kled neovim bulk 5

# Execute code in a terminal
kled neovim code my-terminal --lang python "print('Hello, World!')"

# List available providers
kled neovim providers
```

### Programmatic Usage

The Coder integration can also be used programmatically:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spectrumwebco/agent_runtime/internal/terminal"
	"github.com/spectrumwebco/agent_runtime/pkg/librechat"
)

func main() {
	// Create terminal manager
	libreChatURL := "https://librechat.ai"
	libreChatAPIKey := "your-api-key"
	apiURL := "http://localhost:8080"
	
	terminalManager := terminal.NewManager(apiURL, libreChatURL, libreChatAPIKey)
	
	// Create terminal
	options := map[string]interface{}{
		"api_url": apiURL,
		"cpus":    2,
		"memory":  2048,
	}
	
	term, err := terminalManager.CreateTerminal(context.Background(), "neovim", "my-terminal", options)
	if err != nil {
		log.Fatalf("Failed to create terminal: %v", err)
	}
	
	// Start terminal
	if err := term.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start terminal: %v", err)
	}
	
	// Execute command
	output, err := term.Execute(context.Background(), "ls -la")
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
	
	fmt.Println(output)
}
```

## Kata Container Integration

The Coder integration includes support for Kata Containers, providing enhanced security and isolation for development environments. Kata Containers combine the security advantages of VMs with the performance and density of containers.

### Creating a Kata Container Terminal

```go
// Create a Kata container terminal
options := map[string]interface{}{
	"api_url":         "http://localhost:8080",
	"kata_config_path": "/etc/kata-containers/configuration.toml",
	"kata_runtime_dir": "/var/run/kata-containers",
	"cpus":            2,
	"memory":          2048,
	"debug":           false,
}

terminalManager := terminal.NewManager("http://localhost:8080", "https://librechat.ai", "your-api-key")
term, err := terminalManager.CreateTerminal(context.Background(), "neovim-kata", "my-terminal", options)
if err != nil {
	log.Fatalf("Failed to create terminal: %v", err)
}

if err := term.Start(context.Background()); err != nil {
	log.Fatalf("Failed to start terminal: %v", err)
}
```

## Integration with Other Kled.io Components

The Coder integration works seamlessly with other components of the Kled.io Framework:

1. **Authentication**: Uses the Hydra OAuth2 server for secure authentication
2. **Kubernetes**: Integrates with K9s and Kops for Kubernetes-based environments
3. **CLI**: Follows the unified CLI structure of the Kled.io Framework
4. **Configuration**: Shares the common configuration system

## Security Considerations

The Coder integration includes several security features:

1. **Isolated Environments**: Each development environment is isolated using Kata Containers
2. **Secure Authentication**: Integration with Hydra OAuth2 server
3. **RBAC Support**: Role-based access control for workspaces and templates
4. **Secure Communication**: HTTPS for all API communication

## Future Enhancements

Planned enhancements for the Coder integration include:

1. **Multi-Cloud Support**: Extend support for multiple cloud providers
2. **Custom Templates**: More template options for different development scenarios
3. **Team Collaboration**: Enhanced features for team collaboration
4. **CI/CD Integration**: Tighter integration with CI/CD pipelines

## Conclusion

The Coder integration provides a powerful way for users to create and manage cloud development environments within the Kled.io Framework. By leveraging infrastructure as code and Kata Containers, it offers a secure, flexible, and efficient solution for remote development.
