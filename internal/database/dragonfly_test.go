package database

import (
	"context"
	"testing"
	"time"
	
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestDragonflyAdapter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer s.Close()
	
	connStr := "redis://" + s.Addr()
	
	t.Run("NewDragonflyAdapter", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		assert.NotNil(t, adapter)
		defer adapter.Close()
	})
	
	t.Run("BasicOperations", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		defer adapter.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := adapter.Execute(ctx, "SET testkey testvalue")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		
		result, err = adapter.Query(ctx, "GET testkey")
		assert.NoError(t, err)
		assert.Equal(t, "testvalue", result)
		
		result, err = adapter.Execute(ctx, "DEL testkey")
		assert.NoError(t, err)
		assert.Equal(t, int64(1), result)
		
		result, err = adapter.Query(ctx, "GET testkey")
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})
	
	t.Run("HashOperations", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		defer adapter.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := adapter.Execute(ctx, "HSET testhash field1 value1 field2 value2")
		assert.NoError(t, err)
		assert.Equal(t, int64(2), result)
		
		result, err = adapter.Query(ctx, "HGETALL testhash")
		assert.NoError(t, err)
		hashResult, ok := result.(map[string]string)
		assert.True(t, ok)
		assert.Equal(t, "value1", hashResult["field1"])
		assert.Equal(t, "value2", hashResult["field2"])
	})
	
	t.Run("ListOperations", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		defer adapter.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		result, err := adapter.Execute(ctx, "LPUSH testlist value1 value2 value3")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), result)
		
		result, err = adapter.Query(ctx, "LRANGE testlist 0 -1")
		assert.NoError(t, err)
		listResult, ok := result.([]string)
		assert.True(t, ok)
		assert.Equal(t, 3, len(listResult))
		assert.Equal(t, "value3", listResult[0]) // LPUSH adds to the front
		assert.Equal(t, "value2", listResult[1])
		assert.Equal(t, "value1", listResult[2])
	})
	
	t.Run("ErrorHandling", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		defer adapter.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		_, err = adapter.Query(ctx, "INVALID_COMMAND")
		assert.Error(t, err)
		
		_, err = adapter.Query(ctx, "GET")
		assert.Error(t, err)
	})
	
	t.Run("GetSchema", func(t *testing.T) {
		adapter, err := NewDragonflyAdapter(connStr)
		assert.NoError(t, err)
		defer adapter.Close()
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		_, err = adapter.Execute(ctx, "SET key1 value1")
		assert.NoError(t, err)
		_, err = adapter.Execute(ctx, "HSET hash1 field1 value1")
		assert.NoError(t, err)
		
		schema, err := adapter.GetSchema(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
		
		schemaMap, ok := schema.(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, schemaMap, "info")
		assert.Contains(t, schemaMap, "db_size")
		assert.Contains(t, schemaMap, "key_count")
		assert.Contains(t, schemaMap, "key_types")
	})
}
