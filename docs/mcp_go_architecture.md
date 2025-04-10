# MCP-Go Architecture

## Overview

MCP-Go is a Go implementation of the Model Context Protocol (MCP), which enables AI models to interact with external tools and resources. This document outlines the architecture of MCP-Go and how it can be used to create MCP servers for the Agent Runtime.

## Core Components

### MCP Package

The `mcp` package defines the core types and interfaces for the protocol:

1. **MCPMethod**: Enum of supported MCP methods
2. **Resource**: Represents a data source that can be accessed through MCP
3. **ResourceTemplate**: Template for dynamically generated resources
4. **Tool**: Represents a function that can be executed through MCP
5. **Prompt**: Template for generating text with the model
6. **Content**: Interface for different types of content (text, image, etc.)
7. **JSONRPCMessage**: Interface for JSON-RPC messages
8. **JSONRPCRequest**: Represents a JSON-RPC request
9. **JSONRPCResponse**: Represents a JSON-RPC response
10. **JSONRPCNotification**: Represents a JSON-RPC notification
11. **JSONRPCError**: Represents a JSON-RPC error

### Server Package

The `server` package provides the implementation of the MCP server:

1. **MCPServer**: Main server object that handles requests and responses
2. **ServerOption**: Function that configures an MCPServer
3. **ResourceHandlerFunc**: Function that returns resource contents
4. **ResourceTemplateHandlerFunc**: Function that returns a resource template
5. **PromptHandlerFunc**: Function that handles prompt requests
6. **ToolHandlerFunc**: Function that handles tool calls
7. **NotificationHandlerFunc**: Function that handles notifications
8. **ClientSession**: Interface for client sessions
9. **Hooks**: Provides hooks for server events

## Server Architecture

### Server Creation

The `NewMCPServer` function creates a new MCP server with the specified name, version, and options:

```go
server := server.NewMCPServer(
    "server-name",
    "version",
    server.WithResourceCapabilities(true, true),
    server.WithPromptCapabilities(true),
    server.WithToolCapabilities(true),
    server.WithLogging(),
)
```

### Server Options

Server options configure the behavior of the server:

1. **WithResourceCapabilities**: Enables resource-related capabilities
2. **WithPromptCapabilities**: Enables prompt-related capabilities
3. **WithToolCapabilities**: Enables tool-related capabilities
4. **WithLogging**: Enables logging capabilities
5. **WithPaginationLimit**: Sets the pagination limit for the server
6. **WithInstructions**: Sets the server instructions for the client
7. **WithHooks**: Adds hooks for server events

### Resource Handling

Resources are data sources that can be accessed through MCP:

```go
server.AddResource(mcp.NewResource(
    "resource-uri",
    "Resource Name",
    mcp.WithResourceDescription("Description"),
    mcp.WithMIMEType("text/plain"),
), handleResourceFunc)
```

Resource handlers return resource contents:

```go
func handleResourceFunc(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      request.Params.URI,
            MIMEType: "text/plain",
            Text:     "Resource content",
        },
    }, nil
}
```

### Resource Templates

Resource templates allow for dynamically generated resources:

```go
server.AddResourceTemplate(
    mcp.NewResourceTemplate(
        "template-uri/{param}",
        "Template Name",
    ),
    handleTemplateFunc,
)
```

Template handlers return resource contents based on the template parameters:

```go
func handleTemplateFunc(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    param := request.Params.Arguments["param"].(string)
    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      request.Params.URI,
            MIMEType: "text/plain",
            Text:     fmt.Sprintf("Template content with param: %s", param),
        },
    }, nil
}
```

### Tool Handling

Tools are functions that can be executed through MCP:

```go
server.AddTool(mcp.NewTool(
    "tool-name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param-name",
        mcp.Description("Parameter description"),
        mcp.Required(),
    ),
), handleToolFunc)
```

Tool handlers execute the tool and return the result:

```go
func handleToolFunc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    param := request.Params.Arguments["param-name"].(string)
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Type: "text",
                Text: fmt.Sprintf("Tool result with param: %s", param),
            },
        },
    }, nil
}
```

### Prompt Handling

Prompts are templates for generating text with the model:

```go
server.AddPrompt(mcp.NewPrompt(
    "prompt-name",
    mcp.WithPromptDescription("Prompt description"),
    mcp.WithArgument("param-name",
        mcp.ArgumentDescription("Parameter description"),
        mcp.RequiredArgument(),
    ),
), handlePromptFunc)
```

Prompt handlers generate text based on the prompt and arguments:

```go
func handlePromptFunc(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    param := request.Params.Arguments["param-name"].(string)
    return &mcp.GetPromptResult{
        Description: "Prompt result",
        Messages: []mcp.PromptMessage{
            {
                Role: mcp.RoleUser,
                Content: mcp.TextContent{
                    Type: "text",
                    Text: fmt.Sprintf("Prompt with param: %s", param),
                },
            },
        },
    }, nil
}
```

### Notification Handling

Notifications are asynchronous events that can be sent to the client:

```go
server.AddNotificationHandler("notification-method", handleNotificationFunc)
```

Notification handlers process incoming notifications:

```go
func handleNotificationFunc(ctx context.Context, notification mcp.JSONRPCNotification) {
    // Process notification
}
```

### Server Transport

MCP servers can use different transport mechanisms:

1. **Stdio**: Standard input/output
2. **SSE**: Server-Sent Events over HTTP

```go
// For stdio transport
if err := server.ServeStdio(server); err != nil {
    log.Fatalf("Server error: %v", err)
}

// For SSE transport
sseServer := server.NewSSEServer(server, server.WithBaseURL("http://localhost:8080"))
if err := sseServer.Start(":8080"); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

## MCP Server Types

### Filesystem Server

The Filesystem Server provides access to the local file system through MCP:

```go
type FilesystemServer struct {
    allowedDirs []string
    server      *server.MCPServer
}

func NewFilesystemServer(allowedDirs []string) (*FilesystemServer, error) {
    // Initialize server
    return &FilesystemServer{
        allowedDirs: allowedDirs,
        server:      server.NewMCPServer("filesystem-server", "1.0.0"),
    }, nil
}
```

### LLM Server

The LLM Server provides access to language models through MCP:

```go
type LLMServer struct {
    config LLMConfig
    server *server.MCPServer
}

func NewLLMServer(config LLMConfig) (*LLMServer, error) {
    // Initialize server
    return &LLMServer{
        config: config,
        server: server.NewMCPServer("llm-server", "1.0.0"),
    }, nil
}
```

### Custom Servers

Custom MCP servers can be created to provide specialized capabilities:

```go
type CustomServer struct {
    config CustomConfig
    server *server.MCPServer
}

func NewCustomServer(config CustomConfig) (*CustomServer, error) {
    // Initialize server
    return &CustomServer{
        config: config,
        server: server.NewMCPServer("custom-server", "1.0.0"),
    }, nil
}
```

## Integration with Agent Runtime

### MCP Manager

The Agent Runtime includes an MCP Manager that handles the discovery, connection, and communication with MCP servers:

```go
type Manager struct {
    servers map[string]interface{}
    mutex   sync.RWMutex
}

func NewManager() *Manager {
    return &Manager{
        servers: make(map[string]interface{}),
    }
}

func (m *Manager) RegisterServer(name string, server interface{}) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.servers[name] = server
}

func (m *Manager) CallTool(ctx context.Context, serverName, toolName string, args map[string]interface{}) (*mcp.CallToolResult, error) {
    // Call tool on server
}

func (m *Manager) ReadResource(ctx context.Context, uri string) ([]mcp.ResourceContents, error) {
    // Read resource from server
}
```

### Agent Loop Integration

The Agent Loop can use MCP tools and resources through the MCP Manager:

```go
func (l *Loop) executeToolDiscovery(ctx context.Context) (string, error) {
    // Discover available tools
    tools, err := l.mcpManager.ListTools(ctx)
    if err != nil {
        return "", err
    }
    
    // Select appropriate tools
    selectedTools := selectTools(tools, l.state.Task)
    
    // Store selected tools in state
    l.mutex.Lock()
    l.state.Context["selected_tools"] = selectedTools
    l.mutex.Unlock()
    
    return "planning_phase", nil
}

func (l *Loop) executeExecutionPhase(ctx context.Context) (string, error) {
    // Get selected tool
    l.mutex.RLock()
    toolName := l.state.Context["current_tool"].(string)
    serverName := l.state.Context["current_server"].(string)
    args := l.state.Context["current_args"].(map[string]interface{})
    l.mutex.RUnlock()
    
    // Call tool
    result, err := l.mcpManager.CallTool(ctx, serverName, toolName, args)
    if err != nil {
        return "", err
    }
    
    // Process result
    l.mutex.Lock()
    l.state.Context["current_result"] = result
    l.mutex.Unlock()
    
    return "tool_transition_evaluation", nil
}
```

## Conclusion

MCP-Go provides a powerful framework for creating MCP servers that can be integrated with the Agent Runtime. By understanding the architecture of MCP-Go, developers can create specialized servers that extend the capabilities of the agent.
