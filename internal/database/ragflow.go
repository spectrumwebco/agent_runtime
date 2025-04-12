package database

import (
	"context"
	"database/sql"
	"fmt"
	
	_ "github.com/lib/pq" // PostgreSQL driver
)

type RAGflowAdapter struct {
	connStr string
	db      *sql.DB
}

func NewRAGflowAdapter(connStr string) (*RAGflowAdapter, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open RAGflow connection: %w", err)
	}
	
	return &RAGflowAdapter{
		connStr: connStr,
		db:      db,
	}, nil
}

func (a *RAGflowAdapter) Query(ctx context.Context, query string) (interface{}, error) {
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

func (a *RAGflowAdapter) Execute(ctx context.Context, command string) (interface{}, error) {
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

func (a *RAGflowAdapter) GetSchema(ctx context.Context) (interface{}, error) {
	tableQuery := `
	SELECT table_name 
	FROM information_schema.tables 
	WHERE table_schema = 'public' AND (
		table_name LIKE '%embedding%' OR 
		table_name LIKE '%vector%' OR 
		table_name LIKE '%document%' OR 
		table_name LIKE '%chunk%' OR 
		table_name LIKE '%index%' OR
		table_name LIKE '%rag%'
	)
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
			
			if dataType == "vector" {
				dimensionQuery := `
				SELECT vector_dims($1)
				FROM (SELECT $1::vector FROM ${tableName} LIMIT 1) AS t;
				`
				
				var dimension int
				err := a.db.QueryRowContext(ctx, dimensionQuery, columnName).Scan(&dimension)
				if err == nil {
					columnInfo["vector_dimension"] = dimension
				}
			}
			
			columnList = append(columnList, columnInfo)
		}
		columns.Close()
		
		tableInfo["columns"] = columnList
		tableList = append(tableList, tableInfo)
	}
	
	schema["tables"] = tableList
	
	indexQuery := `
	SELECT indexname, indexdef
	FROM pg_indexes
	WHERE schemaname = 'public' AND indexdef LIKE '%vector%'
	ORDER BY indexname;
	`
	
	indexes, err := a.db.QueryContext(ctx, indexQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query indexes: %w", err)
	}
	defer indexes.Close()
	
	indexList := make([]map[string]interface{}, 0)
	
	for indexes.Next() {
		var indexName, indexDef string
		if err := indexes.Scan(&indexName, &indexDef); err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}
		
		indexInfo := make(map[string]interface{})
		indexInfo["name"] = indexName
		indexInfo["definition"] = indexDef
		
		indexList = append(indexList, indexInfo)
	}
	
	schema["vector_indexes"] = indexList
	
	return schema, nil
}

func (a *RAGflowAdapter) Close() error {
	return a.db.Close()
}

func (a *RAGflowAdapter) GetVectorSearch(ctx context.Context, tableName, vectorColumn, whereClause string, embedding []float64, limit int) (interface{}, error) {
	embeddingStr := "["
	for i, val := range embedding {
		if i > 0 {
			embeddingStr += ","
		}
		embeddingStr += fmt.Sprintf("%f", val)
	}
	embeddingStr += "]"
	
	query := fmt.Sprintf(`
	SELECT *, %s <=> '%s'::vector AS distance
	FROM %s
	WHERE %s
	ORDER BY distance
	LIMIT %d;
	`, vectorColumn, embeddingStr, tableName, whereClause, limit)
	
	return a.Query(ctx, query)
}
