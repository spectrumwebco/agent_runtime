package terminal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBulkTerminalCreation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		switch r.URL.Path {
		case "/start":
			w.Write([]byte(`{"success": true}`))
		case "/stop":
			w.Write([]byte(`{"success": true}`))
		case "/execute":
			w.Write([]byte(`{"success": true, "output": "Command executed"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	manager := NewTerminalManager()

	ctx := context.Background()
	options := map[string]interface{}{
		"api_url": server.URL,
		"id_prefix": "test_",
	}
	
	terminals, err := manager.CreateBulkTerminals(ctx, "neovim", 3, options)
	
	assert.NoError(t, err)
	assert.Equal(t, 3, len(terminals))
	assert.Equal(t, "test_0", terminals[0].ID())
	assert.Equal(t, "test_1", terminals[1].ID())
	assert.Equal(t, "test_2", terminals[2].ID())
	
	for _, terminal := range terminals {
		assert.True(t, terminal.IsRunning())
		assert.Equal(t, "neovim", terminal.GetType())
	}
	
	listedTerminals := manager.ListTerminals()
	assert.Equal(t, 3, len(listedTerminals))
	
	for _, terminal := range terminals {
		err := manager.RemoveTerminal(ctx, terminal.ID())
		assert.NoError(t, err)
	}
	
	assert.Equal(t, 0, len(manager.ListTerminals()))
}

func TestNeovimTerminalBulkExecution(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		switch r.URL.Path {
		case "/start":
			w.Write([]byte(`{"success": true}`))
		case "/stop":
			w.Write([]byte(`{"success": true}`))
		case "/execute":
			w.Write([]byte(`{"success": true, "output": "Command executed"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	terminal1 := NewNeovimTerminal("test_1", server.URL, nil)
	terminal2 := NewNeovimTerminal("test_2", server.URL, nil)
	terminal3 := NewNeovimTerminal("test_3", server.URL, nil)
	
	ctx := context.Background()
	
	err1 := terminal1.Start(ctx)
	err2 := terminal2.Start(ctx)
	err3 := terminal3.Start(ctx)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	
	results := make(map[string]string)
	
	output1, err1 := terminal1.Execute(ctx, "ls -la")
	output2, err2 := terminal2.Execute(ctx, "ls -la")
	output3, err3 := terminal3.Execute(ctx, "ls -la")
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	
	results[terminal1.ID()] = output1
	results[terminal2.ID()] = output2
	results[terminal3.ID()] = output3
	
	assert.Equal(t, 3, len(results))
	assert.Equal(t, "Command executed", results["test_1"])
	assert.Equal(t, "Command executed", results["test_2"])
	assert.Equal(t, "Command executed", results["test_3"])
	
	err1 = terminal1.Stop(ctx)
	err2 = terminal2.Stop(ctx)
	err3 = terminal3.Stop(ctx)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	
	assert.False(t, terminal1.IsRunning())
	assert.False(t, terminal2.IsRunning())
	assert.False(t, terminal3.IsRunning())
}
