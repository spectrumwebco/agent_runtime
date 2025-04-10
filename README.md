# Agent Runtime

<p align="center">
  <img src="docs/assets/logo.png" alt="Agent Runtime Logo" width="200"/>
</p>

<p align="center">
  <a href="https://github.com/spectrumwebco/agent_runtime/actions"><img src="https://github.com/spectrumwebco/agent_runtime/workflows/Go/badge.svg" alt="Build Status"></a>
  <a href="https://pkg.go.dev/github.com/spectrumwebco/agent_runtime"><img src="https://pkg.go.dev/badge/github.com/spectrumwebco/agent_runtime.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/spectrumwebco/agent_runtime"><img src="https://goreportcard.com/badge/github.com/spectrumwebco/agent_runtime" alt="Go Report Card"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
</p>

## Overview

Agent Runtime is a high-performance Go-based backend framework for building autonomous software engineering agents, featuring Sam Sepiol as the lead agent. It leverages the Model Context Protocol (MCP) to create a powerful, extensible system that can integrate with various LLMs and tools.

This framework serves as the foundation for the Agent Runtime ecosystem, providing the core backend logic without frontends. It's designed to be integrated into Kubernetes clusters using kata containers, with native support for TypeScript-based guardrails and mutations through kpolicy.

## Key Features

- **MCP Server Integration**: Native support for the Model Context Protocol, enabling seamless communication between LLMs and external tools
- **Go-Based Performance**: High-throughput, low-latency execution compared to Python-based alternatives
- **Language Interoperability**: Call Python and C++ code directly from Go for compatibility with existing tools and high-performance components
- **Kubernetes Native**: Designed to run in kata containers within a Kubernetes cluster
- **IDE Integration**: Support for JetBrains Toolbox, Windsurf, and VSCode
- **Modular Architecture**: Extensible design with pluggable components
- **Tool Framework**: Comprehensive tool system for agent capabilities
- **Runtime Execution**: Sandboxed environment for safe command execution

## Architecture

Agent Runtime follows a modular architecture with the following key components:

### Core Components

- **Agent Loop**: Manages the agent's execution cycle and state
- **Modules**: Pluggable components for extending functionality
- **MCP Servers**: Specialized servers for different capabilities
- **Runtime Execution**: Sandboxed environment for command execution

### MCP Servers

The framework includes several MCP servers:

- **Filesystem Server**: File operations and management
- **Tool Server**: Execution of agent tools and commands
- **LLM Server**: Integration with various LLM providers
- **Runtime Server**: Execution environment management

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker (for containerized execution)
- Kubernetes cluster (for production deployment)

### Installation

```bash
# Clone the repository
git clone https://github.com/spectrumwebco/agent_runtime.git
cd agent_runtime

# Build the project
go build -o agent_runtime cmd/main.go

# Run the server
./agent_runtime serve
```

### Configuration

Agent Runtime can be configured using environment variables or a configuration file:

```bash
# Using environment variables
export AGENT_RUNTIME_PORT=8080
export AGENT_RUNTIME_LOG_LEVEL=info

# Using configuration file
./agent_runtime serve --config config.yaml
```

Example configuration file:

```yaml
server:
  port: 8080
  host: 0.0.0.0
  
llm:
  provider: openai
  model: gpt-4
  
mcp:
  servers:
    - name: filesystem
      enabled: true
    - name: tools
      enabled: true
```

## Language Interoperability

### Python Integration

Agent Runtime supports calling Python scripts through FFI:

```go
// Example of Python integration
import "github.com/spectrumwebco/agent_runtime/python"

func main() {
    // Initialize Python interpreter
    py := python.NewInterpreter()
    defer py.Close()
    
    // Execute Python code
    result, err := py.Exec("import numpy as np; print(np.array([1, 2, 3]))")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result)
}
```

### C++ Integration

Agent Runtime also supports executing C++ code for high-performance components:

```go
// Example of C++ integration
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/cpp"
)

func main() {
    // Initialize C++ interpreter
    config := cpp.InterpreterConfig{
        Flags: []string{"-std=c++17", "-O2"},
    }
    cppInterp, err := cpp.NewInterpreter(config)
    if err != nil {
        log.Fatal(err)
    }
    defer cppInterp.Close()
    
    // Execute C++ code
    ctx := context.Background()
    code := `
    #include <iostream>
    #include <vector>
    
    int main() {
        std::vector<int> v = {1, 2, 3};
        for (auto i : v) {
            std::cout << i << " ";
        }
        return 0;
    }
    `
    
    result, err := cppInterp.Exec(ctx, code)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(result)
}
```

## Development

### Project Structure

```
agent_runtime/
├── cmd/                    # Command-line applications
├── internal/               # Private application code
│   ├── agent/              # Agent implementation
│   ├── mcp/                # MCP server implementations
│   ├── python/             # Python FFI integration
│   └── runtime/            # Runtime execution environment
├── pkg/                    # Public libraries
│   ├── tools/              # Tool definitions and implementations
│   ├── modules/            # Module system
│   └── config/             # Configuration management
├── api/                    # API definitions
├── web/                    # Web assets and templates
├── scripts/                # Build and deployment scripts
├── docs/                   # Documentation
└── examples/               # Example applications
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Agent Runtime is inspired by and builds upon the work of:

- [SWE-agent](https://github.com/SWE-agent/SWE-agent)
- [SWE-ReX](https://github.com/SWE-agent/SWE-ReX)
- [MCP-Go](https://github.com/mark3labs/mcp-go)
- [MCPHost](https://github.com/mark3labs/mcphost)

## Contact

For questions or support, please open an issue on GitHub or contact the maintainers directly.
