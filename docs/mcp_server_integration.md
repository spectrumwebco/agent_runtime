# MCP Server Integration Guide

## Overview

This document outlines how to create and integrate Model Context Protocol (MCP) servers with the Agent Runtime framework. MCP servers provide specialized capabilities to the agent, allowing it to interact with various systems and services through a standardized protocol.

## MCP Architecture

The Model Context Protocol (MCP) is a standardized communication protocol that enables AI models to interact with external tools and resources. The architecture consists of:

1. **MCP Server**: Implements the protocol and provides specific capabilities
2. **MCP Client**: Consumes the services provided by MCP servers
3. **Resources**: Data sources that can be accessed through the protocol
4. **Tools**: Functions that can be executed through the protocol
5. **Prompts**: Templates for generating text with the model

## Creating MCP Servers

### Core Components

1. **Server Instance**: The main server object that handles requests and responses
2. **Resource Handlers**: Functions that provide access to data sources
3. **Tool Handlers**: Functions that implement specific capabilities
4. **Prompt Handlers**: Functions that generate text based on templates
5. **Notification Handlers**: Functions that handle asynchronous events

### Implementation Steps

1. **Create Server Instance**:
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

2. **Add Resources**:
   ```go
   server.AddResource(mcp.NewResource(
       "resource-uri",
       "Resource Name",
       mcp.WithResourceDescription("Description"),
       mcp.WithMIMEType("text/plain"),
   ), handleResourceFunc)
   ```

3. **Add Resource Templates**:
   ```go
   server.AddResourceTemplate(
       mcp.NewResourceTemplate(
           "template-uri/{param}",
           "Template Name",
       ),
       handleTemplateFunc,
   )
   ```

4. **Add Tools**:
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

5. **Add Prompts**:
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

6. **Add Notification Handlers**:
   ```go
   server.AddNotificationHandler("notification-method", handleNotificationFunc)
   ```

7. **Start the Server**:
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

## Specialized MCP Servers

### Filesystem Server

The Filesystem Server provides access to the local file system through MCP. It includes tools for:

- Reading and writing files
- Listing directories
- Creating directories
- Moving files
- Searching for files
- Getting file information

Implementation example:
```go
server := NewFilesystemServer([]string{"/allowed/path"})
if err := server.Serve(); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

### LLM Server

The LLM Server provides access to language models through MCP. It includes tools for:

- Generating text
- Completing prompts
- Answering questions
- Summarizing text

Implementation example:
```go
server := NewLLMServer(llmConfig)
if err := server.Serve(); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

### Custom MCP Servers

Custom MCP servers can be created to provide specialized capabilities for the agent. Examples include:

- Database servers
- API servers
- Visualization servers
- Memory servers
- Search servers

## Integrating MCP Servers with Agent Runtime

### MCP Manager

The Agent Runtime includes an MCP Manager that handles the discovery, connection, and communication with MCP servers. The manager:

1. Discovers available MCP servers
2. Establishes connections to servers
3. Routes requests to appropriate servers
4. Handles server responses
5. Manages server lifecycle

Implementation:
```go
manager := mcp.NewManager()
manager.RegisterServer("filesystem", filesystemServer)
manager.RegisterServer("llm", llmServer)
```

### Using MCP Tools in Agent Loop

The Agent Loop can use MCP tools through the MCP Manager:

```go
result, err := manager.CallTool(ctx, "filesystem", "read_file", map[string]interface{}{
    "path": "/path/to/file",
})
if err != nil {
    // Handle error
}
// Process result
```

### Using MCP Resources in Agent Loop

The Agent Loop can access MCP resources through the MCP Manager:

```go
contents, err := manager.ReadResource(ctx, "file:///path/to/file")
if err != nil {
    // Handle error
}
// Process contents
```

## Integration with kpolicy for Kubernetes

The Agent Runtime can be integrated with kpolicy to create guardrails and mutations for Kubernetes clusters:

1. **MCP Server for Kubernetes**: Create an MCP server that provides access to Kubernetes resources and operations
2. **kpolicy Integration**: Use kpolicy to define guardrails and mutations for Kubernetes resources
3. **TypeScript Integration**: Use TypeScript to define policies and rules for the agent

Implementation example:
```go
// Create Kubernetes MCP server
k8sServer := NewKubernetesServer(k8sConfig)

// Register with MCP Manager
manager.RegisterServer("kubernetes", k8sServer)

// Use in agent loop
result, err := manager.CallTool(ctx, "kubernetes", "apply", map[string]interface{}{
    "manifest": manifest,
    "namespace": namespace,
})
```

## Best Practices

1. **Security**: Validate all inputs and restrict access to sensitive resources
2. **Error Handling**: Provide clear error messages and handle errors gracefully
3. **Performance**: Optimize for performance, especially for frequently used tools
4. **Logging**: Log all operations for debugging and auditing
5. **Testing**: Write comprehensive tests for all MCP servers and tools
6. **Documentation**: Document all resources, tools, and prompts
7. **Versioning**: Use semantic versioning for MCP servers
8. **Compatibility**: Ensure compatibility with different MCP clients

## Conclusion

MCP servers provide a powerful way to extend the capabilities of the Agent Runtime. By creating specialized servers for different domains, the agent can interact with a wide range of systems and services through a standardized protocol.
