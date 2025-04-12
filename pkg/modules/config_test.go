package modules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

func TestModuleConfigInterchangeability(t *testing.T) {
	tempDir := t.TempDir()
	goConfigsDir := filepath.Join(tempDir, "go_configs")
	pkgDir := filepath.Join(tempDir, "pkg")
	pkgPromptsDir := filepath.Join(pkgDir, "prompts")
	pkgToolsDir := filepath.Join(pkgDir, "tools")
	
	if err := os.MkdirAll(goConfigsDir, 0755); err != nil {
		t.Fatalf("Failed to create go_configs directory: %v", err)
	}
	if err := os.MkdirAll(pkgPromptsDir, 0755); err != nil {
		t.Fatalf("Failed to create pkg/prompts directory: %v", err)
	}
	if err := os.MkdirAll(pkgToolsDir, 0755); err != nil {
		t.Fatalf("Failed to create pkg/tools directory: %v", err)
	}
	
	goConfigsPromptsContent := `# Test prompts from go_configs
You are samsepi0l, an autonomous software engineering agent.
`
	pkgPromptsContent := `# Test prompts from pkg
You are samsepi0l, an autonomous software engineering agent.
`
	
	goConfigsToolsContent := `{
  "tools": [
    {
      "name": "test_tool",
      "description": "A test tool from go_configs",
      "parameters": {}
    }
  ]
}`
	
	pkgToolsContent := `{
  "tools": [
    {
      "name": "test_tool",
      "description": "A test tool from pkg",
      "parameters": {}
    }
  ]
}`
	
	if err := os.WriteFile(filepath.Join(goConfigsDir, "prompts.txt"), []byte(goConfigsPromptsContent), 0644); err != nil {
		t.Fatalf("Failed to write go_configs prompts.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkgPromptsDir, "prompts.txt"), []byte(pkgPromptsContent), 0644); err != nil {
		t.Fatalf("Failed to write pkg prompts.txt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(goConfigsDir, "tools.json"), []byte(goConfigsToolsContent), 0644); err != nil {
		t.Fatalf("Failed to write go_configs tools.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkgToolsDir, "tools.json"), []byte(pkgToolsContent), 0644); err != nil {
		t.Fatalf("Failed to write pkg tools.json: %v", err)
	}
	
	originalFindRepoRoot := config.FindRepoRoot
	defer func() { config.FindRepoRoot = originalFindRepoRoot }()
	config.FindRepoRoot = func() (string, error) {
		return tempDir, nil
	}
	
	t.Run("Load from go_configs", func(t *testing.T) {
		adapter, err := config.NewConfigAdapter(config.GoConfigsSource)
		if err != nil {
			t.Fatalf("Failed to create config adapter: %v", err)
		}
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts from go_configs: %v", err)
		}
		
		if !strings.Contains(prompts, "samsepi0l") {
			t.Errorf("Expected prompts to contain 'samsepi0l', got: %s", prompts)
		}
		
		toolsConfig, err := adapter.LoadTools()
		if err != nil {
			t.Fatalf("Failed to load tools from go_configs: %v", err)
		}
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if !ok {
			t.Fatalf("Invalid tools format")
		}
		
		if len(toolsArray) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(toolsArray))
		}
	})
	
	t.Run("Load from pkg", func(t *testing.T) {
		adapter, err := config.NewConfigAdapter(config.PkgSource)
		if err != nil {
			t.Fatalf("Failed to create config adapter: %v", err)
		}
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts from pkg: %v", err)
		}
		
		if !strings.Contains(prompts, "samsepi0l") {
			t.Errorf("Expected prompts to contain 'samsepi0l', got: %s", prompts)
		}
		
		toolsConfig, err := adapter.LoadTools()
		if err != nil {
			t.Fatalf("Failed to load tools from pkg: %v", err)
		}
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if !ok {
			t.Fatalf("Invalid tools format")
		}
		
		if len(toolsArray) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(toolsArray))
		}
	})
	
	t.Run("Test fallback mechanism", func(t *testing.T) {
		if err := os.Remove(filepath.Join(goConfigsDir, "prompts.txt")); err != nil {
			t.Fatalf("Failed to remove go_configs prompts.txt: %v", err)
		}
		
		adapter, err := config.NewConfigAdapter(config.GoConfigsSource)
		if err != nil {
			t.Fatalf("Failed to create config adapter: %v", err)
		}
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts with fallback: %v", err)
		}
		
		if !strings.Contains(prompts, "samsepi0l") {
			t.Errorf("Expected prompts to contain 'samsepi0l', got: %s", prompts)
		}
	})
	
	t.Run("Test module registry with tools", func(t *testing.T) {
		registry := NewRegistry()
		
		module := &BaseModule{
			name:        "test_module",
			description: "Test module for configuration interchangeability",
		}
		
		registry.Register(module)
		
		modules := registry.List()
		if len(modules) != 1 {
			t.Errorf("Expected 1 module, got %d", len(modules))
		}
		
		adapter, _ := config.NewConfigAdapter(config.GoConfigsSource)
		toolsConfig, _ := adapter.LoadTools()
		
		toolRegistry := &tools.ToolRegistry{}
		err := tools.LoadToolsFromConfig(toolRegistry, toolsConfig)
		if err != nil {
			t.Fatalf("Failed to load tools from config: %v", err)
		}
		
		commands := toolRegistry.ListCommands()
		if len(commands) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(commands))
		}
		
		adapter, _ = config.NewConfigAdapter(config.PkgSource)
		toolsConfig, _ = adapter.LoadTools()
		
		toolRegistry = &tools.ToolRegistry{}
		err = tools.LoadToolsFromConfig(toolRegistry, toolsConfig)
		if err != nil {
			t.Fatalf("Failed to load tools from config: %v", err)
		}
		
		commands = toolRegistry.ListCommands()
		if len(commands) != 1 {
			t.Errorf("Expected 1 tool, got %d", len(commands))
		}
	})
}
