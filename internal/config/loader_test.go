package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	goConfigsDir := filepath.Join(tempDir, "go_configs")
	pkgDir := filepath.Join(tempDir, "pkg")
	pkgPromptsDir := filepath.Join(pkgDir, "prompts")
	pkgToolsDir := filepath.Join(pkgDir, "tools")
	
	os.MkdirAll(goConfigsDir, 0755)
	os.MkdirAll(pkgPromptsDir, 0755)
	os.MkdirAll(pkgToolsDir, 0755)
	
	goConfigsPromptsContent := "# Test prompts from go_configs\nYou are samsepi0l, a software engineer."
	pkgPromptsContent := "# Test prompts from pkg\nYou are samsepi0l, a software engineer."
	
	goConfigsToolsContent := `{"tools": [{"name": "test_tool", "description": "A test tool"}]}`
	pkgToolsContent := `{"tools": [{"name": "test_tool", "description": "A test tool from pkg"}]}`
	
	os.WriteFile(filepath.Join(goConfigsDir, "prompts.txt"), []byte(goConfigsPromptsContent), 0644)
	os.WriteFile(filepath.Join(pkgPromptsDir, "prompts.txt"), []byte(pkgPromptsContent), 0644)
	os.WriteFile(filepath.Join(goConfigsDir, "tools.json"), []byte(goConfigsToolsContent), 0644)
	os.WriteFile(filepath.Join(pkgToolsDir, "tools.json"), []byte(pkgToolsContent), 0644)
	
	originalFindRepoRoot := findRepoRoot
	defer func() { findRepoRoot = originalFindRepoRoot }()
	findRepoRoot = func() string { return tempDir }
	
	t.Run("Load prompts from go_configs", func(t *testing.T) {
		data, err := LoadConfig(PromptsConfig, GoConfigsLocation, PkgLocation)
		if err != nil {
			t.Fatalf("Failed to load prompts from go_configs: %v", err)
		}
		
		if string(data) != goConfigsPromptsContent {
			t.Errorf("Expected %q, got %q", goConfigsPromptsContent, string(data))
		}
	})
	
	t.Run("Load prompts from pkg", func(t *testing.T) {
		data, err := LoadConfig(PromptsConfig, PkgLocation, GoConfigsLocation)
		if err != nil {
			t.Fatalf("Failed to load prompts from pkg: %v", err)
		}
		
		if string(data) != pkgPromptsContent {
			t.Errorf("Expected %q, got %q", pkgPromptsContent, string(data))
		}
	})
	
	t.Run("Load tools from go_configs", func(t *testing.T) {
		data, err := LoadConfig(ToolsConfig, GoConfigsLocation, PkgLocation)
		if err != nil {
			t.Fatalf("Failed to load tools from go_configs: %v", err)
		}
		
		if string(data) != goConfigsToolsContent {
			t.Errorf("Expected %q, got %q", goConfigsToolsContent, string(data))
		}
	})
	
	t.Run("Load tools from pkg", func(t *testing.T) {
		data, err := LoadConfig(ToolsConfig, PkgLocation, GoConfigsLocation)
		if err != nil {
			t.Fatalf("Failed to load tools from pkg: %v", err)
		}
		
		if string(data) != pkgToolsContent {
			t.Errorf("Expected %q, got %q", pkgToolsContent, string(data))
		}
	})
	
	t.Run("Test fallback mechanism", func(t *testing.T) {
		os.Remove(filepath.Join(goConfigsDir, "prompts.txt"))
		
		data, err := LoadConfig(PromptsConfig, GoConfigsLocation, PkgLocation)
		if err != nil {
			t.Fatalf("Failed to load prompts with fallback: %v", err)
		}
		
		if string(data) != pkgPromptsContent {
			t.Errorf("Expected fallback content %q, got %q", pkgPromptsContent, string(data))
		}
	})
}

func TestLoadPrompts(t *testing.T) {
	tempDir := t.TempDir()
	goConfigsDir := filepath.Join(tempDir, "go_configs")
	pkgDir := filepath.Join(tempDir, "pkg", "prompts")
	
	os.MkdirAll(goConfigsDir, 0755)
	os.MkdirAll(pkgDir, 0755)
	
	goConfigsPromptsContent := "# Test prompts from go_configs\nYou are samsepi0l, a software engineer."
	pkgPromptsContent := "# Test prompts from pkg\nYou are samsepi0l, a software engineer."
	
	os.WriteFile(filepath.Join(goConfigsDir, "prompts.txt"), []byte(goConfigsPromptsContent), 0644)
	os.WriteFile(filepath.Join(pkgDir, "prompts.txt"), []byte(pkgPromptsContent), 0644)
	
	originalFindRepoRoot := findRepoRoot
	defer func() { findRepoRoot = originalFindRepoRoot }()
	findRepoRoot = func() string { return tempDir }
	
	t.Run("LoadPrompts from go_configs", func(t *testing.T) {
		prompts, err := LoadPrompts(GoConfigsLocation)
		if err != nil {
			t.Fatalf("Failed to load prompts: %v", err)
		}
		
		if prompts != goConfigsPromptsContent {
			t.Errorf("Expected %q, got %q", goConfigsPromptsContent, prompts)
		}
	})
	
	t.Run("LoadPrompts from pkg", func(t *testing.T) {
		prompts, err := LoadPrompts(PkgLocation)
		if err != nil {
			t.Fatalf("Failed to load prompts: %v", err)
		}
		
		if prompts != pkgPromptsContent {
			t.Errorf("Expected %q, got %q", pkgPromptsContent, prompts)
		}
	})
}

func TestLoadTools(t *testing.T) {
	tempDir := t.TempDir()
	goConfigsDir := filepath.Join(tempDir, "go_configs")
	pkgDir := filepath.Join(tempDir, "pkg", "tools")
	
	os.MkdirAll(goConfigsDir, 0755)
	os.MkdirAll(pkgDir, 0755)
	
	goConfigsToolsContent := `{"tools": [{"name": "test_tool", "description": "A test tool"}]}`
	pkgToolsContent := `{"tools": [{"name": "test_tool", "description": "A test tool from pkg"}]}`
	
	os.WriteFile(filepath.Join(goConfigsDir, "tools.json"), []byte(goConfigsToolsContent), 0644)
	os.WriteFile(filepath.Join(pkgDir, "tools.json"), []byte(pkgToolsContent), 0644)
	
	originalFindRepoRoot := findRepoRoot
	defer func() { findRepoRoot = originalFindRepoRoot }()
	findRepoRoot = func() string { return tempDir }
	
	t.Run("LoadTools from go_configs", func(t *testing.T) {
		tools, err := LoadTools(GoConfigsLocation)
		if err != nil {
			t.Fatalf("Failed to load tools: %v", err)
		}
		
		toolsArray, ok := tools["tools"].([]interface{})
		if !ok || len(toolsArray) == 0 {
			t.Fatalf("Invalid tools format or empty tools array")
		}
		
		firstTool, ok := toolsArray[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid tool format")
		}
		
		if firstTool["name"] != "test_tool" {
			t.Errorf("Expected tool name 'test_tool', got %v", firstTool["name"])
		}
	})
	
	t.Run("LoadTools from pkg", func(t *testing.T) {
		tools, err := LoadTools(PkgLocation)
		if err != nil {
			t.Fatalf("Failed to load tools: %v", err)
		}
		
		toolsArray, ok := tools["tools"].([]interface{})
		if !ok || len(toolsArray) == 0 {
			t.Fatalf("Invalid tools format or empty tools array")
		}
		
		firstTool, ok := toolsArray[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid tool format")
		}
		
		if firstTool["description"] != "A test tool from pkg" {
			t.Errorf("Expected tool description 'A test tool from pkg', got %v", firstTool["description"])
		}
	})
}
