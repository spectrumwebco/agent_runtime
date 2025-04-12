package mcp

import (
	"context"
	"fmt"
	
	"github.com/mark3labs/mcp-go/server"
	"github.com/spectrumwebco/agent_runtime/internal/database"
)

func RegisterDatabaseTools(mcpServer *server.MCPServer) error {
	mcpServer.RegisterTool("db_query", "Executes a SQL query using natural language", []server.ToolParameter{
		{
			Name:        "question",
			Type:        "string",
			Description: "Natural language question to ask the database",
			Required:    true,
		},
		{
			Name:        "database_type",
			Type:        "string",
			Description: "Type of database (supabase, dragonfly, ragflow, rocketmq)",
			Required:    true,
		},
		{
			Name:        "connection_string",
			Type:        "string",
			Description: "Connection string for the database",
			Required:    true,
		},
	}, handleDatabaseQuery)
	
	mcpServer.RegisterTool("db_execute", "Executes a SQL command using natural language", []server.ToolParameter{
		{
			Name:        "command",
			Type:        "string",
			Description: "Natural language command to execute on the database",
			Required:    true,
		},
		{
			Name:        "database_type",
			Type:        "string",
			Description: "Type of database (supabase, dragonfly, ragflow, rocketmq)",
			Required:    true,
		},
		{
			Name:        "connection_string",
			Type:        "string",
			Description: "Connection string for the database",
			Required:    true,
		},
	}, handleDatabaseExecute)
	
	mcpServer.RegisterTool("db_connect", "Connects to a database", []server.ToolParameter{
		{
			Name:        "database_type",
			Type:        "string",
			Description: "Type of database (supabase, dragonfly, ragflow, rocketmq)",
			Required:    true,
		},
		{
			Name:        "connection_string",
			Type:        "string",
			Description: "Connection string for the database",
			Required:    true,
		},
	}, handleDatabaseConnect)
	
	mcpServer.RegisterTool("db_schema", "Gets the schema of a database", []server.ToolParameter{
		{
			Name:        "database_type",
			Type:        "string",
			Description: "Type of database (supabase, dragonfly, ragflow, rocketmq)",
			Required:    true,
		},
		{
			Name:        "connection_string",
			Type:        "string",
			Description: "Connection string for the database",
			Required:    true,
		},
	}, handleDatabaseSchema)
	
	return nil
}

func handleDatabaseQuery(params map[string]interface{}) (interface{}, error) {
	question, ok := params["question"].(string)
	if !ok {
		return nil, fmt.Errorf("question parameter is required")
	}
	
	dbType, ok := params["database_type"].(string)
	if !ok {
		return nil, fmt.Errorf("database_type parameter is required")
	}
	
	connStr, ok := params["connection_string"].(string)
	if !ok {
		return nil, fmt.Errorf("connection_string parameter is required")
	}
	
	sql, err := database.NaturalLanguageToSQL(context.Background(), question, dbType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert natural language to SQL: %w", err)
	}
	
	result, err := database.ExecuteQuery(context.Background(), dbType, connStr, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	
	return result, nil
}

func handleDatabaseExecute(params map[string]interface{}) (interface{}, error) {
	command, ok := params["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}
	
	dbType, ok := params["database_type"].(string)
	if !ok {
		return nil, fmt.Errorf("database_type parameter is required")
	}
	
	connStr, ok := params["connection_string"].(string)
	if !ok {
		return nil, fmt.Errorf("connection_string parameter is required")
	}
	
	sql, err := database.NaturalLanguageToSQL(context.Background(), command, dbType)
	if err != nil {
		return nil, fmt.Errorf("failed to convert natural language to SQL: %w", err)
	}
	
	result, err := database.ExecuteCommand(context.Background(), dbType, connStr, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}
	
	return result, nil
}

func handleDatabaseConnect(params map[string]interface{}) (interface{}, error) {
	dbType, ok := params["database_type"].(string)
	if !ok {
		return nil, fmt.Errorf("database_type parameter is required")
	}
	
	connStr, ok := params["connection_string"].(string)
	if !ok {
		return nil, fmt.Errorf("connection_string parameter is required")
	}
	
	err := database.TestConnection(context.Background(), dbType, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return map[string]interface{}{
		"status":  "connected",
		"message": fmt.Sprintf("Successfully connected to %s database", dbType),
	}, nil
}

func handleDatabaseSchema(params map[string]interface{}) (interface{}, error) {
	dbType, ok := params["database_type"].(string)
	if !ok {
		return nil, fmt.Errorf("database_type parameter is required")
	}
	
	connStr, ok := params["connection_string"].(string)
	if !ok {
		return nil, fmt.Errorf("connection_string parameter is required")
	}
	
	schema, err := database.GetSchema(context.Background(), dbType, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get database schema: %w", err)
	}
	
	return map[string]interface{}{
		"schema": schema,
	}, nil
}
