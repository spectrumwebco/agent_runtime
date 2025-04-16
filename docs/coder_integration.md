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

## Directory Structure

```
agent_runtime/
├── internal/
│   └── devenv/
│       └── coder/
│           ├── client.go         # Coder API client
│           ├── provisioner.go    # Workspace provisioning
│           └── terraform.go      # Terraform provider with Kata support
├── pkg/
│   └── modules/
│       └── coder/
│           └── module.go         # Framework module integration
└── cmd/
    └── cli/
        └── commands/
            └── coder.go          # CLI commands
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
kled coder workspace list

# Create a new workspace
kled coder workspace create my-workspace --template template-123

# Get workspace details
kled coder workspace get ws-123

# Start a workspace
kled coder workspace start ws-123

# Stop a workspace
kled coder workspace stop ws-123

# Delete a workspace
kled coder workspace delete ws-123

# List templates
kled coder template list

# Create a Kata container template
kled coder kata template my-template --image ubuntu:20.04 --cpu 2 --memory 4096 --disk 10
```

### Programmatic Usage

The Coder integration can also be used programmatically:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spectrumwebco/agent_runtime/internal/devenv/coder"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create client
	client, err := coder.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List workspaces
	workspaces, err := client.ListWorkspaces(context.Background())
	if err != nil {
		log.Fatalf("Failed to list workspaces: %v", err)
	}

	// Print workspaces
	for _, workspace := range workspaces {
		fmt.Printf("Workspace: %s (Status: %s)\n", workspace.Name, workspace.Status)
	}
}
```

## Kata Container Integration

The Coder integration includes support for Kata Containers, providing enhanced security and isolation for development environments. Kata Containers combine the security advantages of VMs with the performance and density of containers.

### Creating a Kata Container Template

```go
// Create a Kata container template
terraformProvider := coder.NewTerraformProvider(cfg)
config := &coder.KataContainerConfig{
    Image:  "ubuntu:20.04",
    CPU:    2,
    Memory: 4096,
    Disk:   10,
}
templateDir, err := terraformProvider.GenerateKataContainerTemplate(
    context.Background(), 
    "my-template", 
    "A development environment with Kata Containers", 
    config,
)
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
