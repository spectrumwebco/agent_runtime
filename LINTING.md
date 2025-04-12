# Agent Runtime Linting Issues

This document summarizes the linting issues identified in the agent_runtime codebase.

## Interface Redeclarations

1. **Module Interface Redeclaration**
   - Files: `pkg/modules/base.go` and `pkg/modules/module.go`
   - Issue: The `Module` interface is defined in both files with different method signatures
   - Resolution: Consolidate into a single interface definition

2. **Runtime Interface Redeclaration**
   - Files: `internal/agent_framework/abstract.go` and `internal/agent_framework/runtime.go`
   - Issue: The `Runtime` interface is defined in both files
   - Resolution: Consolidate into a single interface definition

## Undefined Methods and Types

1. **MCP Server Methods**
   - Files: `internal/mcp/filesystem_resources.go` and `internal/mcp/filesystem_tools.go`
   - Issues:
     - `mcpServer.RegisterResourceHandler` undefined
     - `mcpServer.RegisterTool` undefined
     - `server.ToolParameter` undefined

2. **Tools Package**
   - Files: `internal/tools/tools.go` and `internal/tools/specialized.go`
   - Issues:
     - `tools.Logger` undefined
     - `tools.NewBaseTool` undefined
     - `tools.DataProcessingTools` undefined
     - `tools.SpecializedTool` undefined
     - `tools.ToolParameter` undefined
     - `registry.RegisterTool` undefined

3. **K8s Package**
   - File: `k8s/k8s_handlers.go`
   - Issue: `server.NewTool` undefined

## Parser Package Issues

1. **Test Failures**
   - File: `internal/parser/factory_test.go`
   - Issue: `assert.TypeOf` undefined

## Next Steps

1. Consolidate interface definitions to eliminate redeclarations
2. Implement missing methods and types
3. Fix test dependencies
4. Add proper documentation comments to exported types and methods
5. Ensure consistent naming and method signatures across packages
