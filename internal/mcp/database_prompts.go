package mcp

import (
	"context"
	"fmt"
	
	"github.com/mark3labs/mcp-go/server"
)

func RegisterDatabasePrompts(mcpServer *server.MCPServer) error {
	mcpServer.RegisterPrompt("db_query_prompt", server.PromptWithDescription("A prompt for converting natural language to SQL"), []server.PromptParameter{
		{
			Name:        "question",
			Type:        "string",
			Description: "Natural language question to convert to SQL",
			Required:    true,
		},
		{
			Name:        "database_type",
			Type:        "string",
			Description: "Type of database (supabase, dragonfly, ragflow, rocketmq)",
			Required:    true,
		},
		{
			Name:        "schema",
			Type:        "string",
			Description: "Database schema information",
			Required:    false,
		},
	}, handleDatabaseQueryPrompt)
	
	return nil
}

func handleDatabaseQueryPrompt(ctx context.Context, req server.GetPromptRequest) (*server.GetPromptResult, error) {
	question, ok := req.Params.Args["question"].(string)
	if !ok {
		return nil, fmt.Errorf("question parameter is required")
	}
	
	dbType, ok := req.Params.Args["database_type"].(string)
	if !ok {
		return nil, fmt.Errorf("database_type parameter is required")
	}
	
	schema, _ := req.Params.Args["schema"].(string)
	
	systemPrompt := fmt.Sprintf(`You are an expert SQL writer for %s databases. 
Given a natural language question, your task is to generate the correct SQL query to answer it.
If schema information is provided, use it to construct your query.

Database Type: %s
Schema Information: %s

Natural Language Question: %s

Respond with ONLY the SQL query, nothing else.`, dbType, dbType, schema, question)
	
	return &server.GetPromptResult{
		Messages: []server.PromptMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
		},
	}, nil
}
