package database

import (
	"context"
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq" // PostgreSQL driver
)

type SupabaseAdapter struct {
	connStr string
	db      *sql.DB
}

func NewSupabaseAdapter(connStr string) (*SupabaseAdapter, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open Supabase connection: %w", err)
	}
	
	return &SupabaseAdapter{
		connStr: connStr,
		db:      db,
	}, nil
}

func (a *SupabaseAdapter) Query(ctx context.Context, query string) (interface{}, error) {
	rows, err := a.db.QueryContext(ctx, query)
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

func (a *SupabaseAdapter) Execute(ctx context.Context, command string) (interface{}, error) {
	result, err := a.db.ExecContext(ctx, command)
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

func (a *SupabaseAdapter) GetSchema(ctx context.Context) (interface{}, error) {
	tableQuery := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema = 'public'
	ORDER BY table_name;
	`
	
	tables, err := a.db.QueryContext(ctx, tableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer tables.Close()
	
	schema := make(map[string]interface{})
	tableList := make([]map[string]interface{}, 0)
	
	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		
		tableInfo := make(map[string]interface{})
		tableInfo["name"] = tableName
		
		columnQuery := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position;
		`
		
		columns, err := a.db.QueryContext(ctx, columnQuery, tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to query columns: %w", err)
		}
		
		columnList := make([]map[string]interface{}, 0)
		
		for columns.Next() {
			var columnName, dataType, isNullable string
			if err := columns.Scan(&columnName, &dataType, &isNullable); err != nil {
				columns.Close()
				return nil, fmt.Errorf("failed to scan column: %w", err)
			}
			
			columnInfo := make(map[string]interface{})
			columnInfo["name"] = columnName
			columnInfo["type"] = dataType
			columnInfo["nullable"] = isNullable == "YES"
			
			columnList = append(columnList, columnInfo)
		}
		columns.Close()
		
		tableInfo["columns"] = columnList
		tableList = append(tableList, tableInfo)
	}
	
	schema["tables"] = tableList
	return schema, nil
}

func (a *SupabaseAdapter) Close() error {
	return a.db.Close()
}
