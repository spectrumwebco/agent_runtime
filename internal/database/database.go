package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	
	"github.com/spectrumwebco/agent_runtime/internal/ai"
)

var (
	connPool     = make(map[string]*sql.DB)
	connPoolLock sync.Mutex
)

func NaturalLanguageToSQL(ctx context.Context, question, dbType string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert SQL writer for %s databases. 
Given a natural language question, your task is to generate the correct SQL query to answer it.

Natural Language Question: %s

Respond with ONLY the SQL query, nothing else.`, dbType, question)
	
	response, err := ai.CompletionWithLLM(ctx, prompt, "llama-4")
	if err != nil {
		return "", fmt.Errorf("failed to generate SQL: %w", err)
	}
	
	return response, nil
}

func ExecuteQuery(ctx context.Context, dbType, connStr, query string) (interface{}, error) {
	db, err := getConnection(dbType, connStr)
	if err != nil {
		return nil, err
	}
	
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	
	var results []map[string]interface{}
	
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		
		row := make(map[string]interface{})
		
		for i, col := range columns {
			val := values[i]
			
			if val == nil {
				row[col] = nil
			} else {
				if b, ok := val.([]byte); ok {
					row[col] = string(b)
				} else {
					row[col] = val
				}
			}
		}
		
		results = append(results, row)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	
	return results, nil
}

func ExecuteCommand(ctx context.Context, dbType, connStr, command string) (interface{}, error) {
	db, err := getConnection(dbType, connStr)
	if err != nil {
		return nil, err
	}
	
	result, err := db.ExecContext(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	return map[string]interface{}{
		"rows_affected": rowsAffected,
	}, nil
}

func TestConnection(ctx context.Context, dbType, connStr string) error {
	db, err := getConnection(dbType, connStr)
	if err != nil {
		return err
	}
	
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	return nil
}

func GetSchema(ctx context.Context, dbType, connStr string) (interface{}, error) {
	db, err := getConnection(dbType, connStr)
	if err != nil {
		return nil, err
	}
	
	var schema interface{}
	
	switch dbType {
	case "supabase", "postgres", "ragflow":
		query := `
			SELECT 
				table_schema,
				table_name,
				column_name,
				data_type,
				is_nullable
			FROM 
				information_schema.columns
			WHERE 
				table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY 
				table_schema, table_name, ordinal_position;
		`
		schema, err = ExecuteQuery(ctx, dbType, connStr, query)
		
	case "dragonfly", "redis":
		schema = map[string]interface{}{
			"message": "Dragonfly/Redis doesn't have a traditional schema. Use key patterns to explore data.",
		}
		
	case "rocketmq":
		schema = map[string]interface{}{
			"message": "RocketMQ schema consists of topics and consumer groups.",
		}
		
	default:
		return nil, fmt.Errorf("unsupported database type for schema retrieval: %s", dbType)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schema: %w", err)
	}
	
	return schema, nil
}

func getConnection(dbType, connStr string) (*sql.DB, error) {
	key := fmt.Sprintf("%s:%s", dbType, connStr)
	
	connPoolLock.Lock()
	defer connPoolLock.Unlock()
	
	if db, ok := connPool[key]; ok {
		return db, nil
	}
	
	var db *sql.DB
	var err error
	
	switch dbType {
	case "supabase", "postgres", "ragflow":
		db, err = sql.Open("postgres", connStr)
	case "dragonfly", "redis":
		db, err = sql.Open("redis", connStr)
	case "rocketmq":
		return nil, fmt.Errorf("rocketmq requires special handling")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	connPool[key] = db
	
	return db, nil
}

func CloseConnections() {
	connPoolLock.Lock()
	defer connPoolLock.Unlock()
	
	for _, db := range connPool {
		db.Close()
	}
	
	connPool = make(map[string]*sql.DB)
}
