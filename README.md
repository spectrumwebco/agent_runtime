# Agent Runtime

Agent Runtime is a Go implementation of SWE-Agent and SWE-ReX frameworks, providing a high-performance agent runtime with Python FFI capabilities.

## Project Structure

The project follows standard Go project layout:

- `cmd/` - Main applications
- `internal/` - Private application and library code
- `pkg/` - Library code that's ok to use by external applications
- `configs/` - Configuration files
- `python_agent/` - Python implementation of the agent

## Python Agent

The Python agent is located in the `python_agent/` directory and includes:

- `agent/` - Core agent implementation
- `tools/` - Tools used by the agent
- `tests/` - Tests for the agent
- `agent_config/` - Configuration files
- `agent_framework/` - Agent framework implementation (formerly SWE-ReX)

## Getting Started

See the documentation in the `docs/` directory for more information.

## Features

- High-performance Go implementation of the agent runtime
- Python FFI for tool execution
- Support for multiple MCP servers
- Integration with Gitee and GitHub repositories
- Kubernetes and Kata container support
