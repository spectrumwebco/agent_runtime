package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigAdapter(t *testing.T) {
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
	defer func() {
		findRepoRoot = originalFindRepoRoot
	}()
	
	findRepoRoot = func() (string, error) {
		return tempDir, nil
	}
	
	t.Run("Load from go_configs", func(t *testing.T) {
		adapter, err := NewConfigAdapter(GoConfigsSource)
		if err != nil {
			t.Fatalf("Failed to create adapter: %v", err)
		}
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts: %v", err)
		}
		
		if prompts != goConfigsPromptsContent {
			t.Errorf("Expected %q, got %q", goConfigsPromptsContent, prompts)
		}
		
		tools, err := adapter.LoadTools()
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
		
		if firstTool["description"] != "A test tool" {
			t.Errorf("Expected tool description 'A test tool', got %v", firstTool["description"])
		}
	})
	
	t.Run("Load from pkg", func(t *testing.T) {
		adapter, err := NewConfigAdapter(PkgSource)
		if err != nil {
			t.Fatalf("Failed to create adapter: %v", err)
		}
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts: %v", err)
		}
		
		if prompts != pkgPromptsContent {
			t.Errorf("Expected %q, got %q", pkgPromptsContent, prompts)
		}
		
		tools, err := adapter.LoadTools()
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
	
	t.Run("Test fallback mechanism", func(t *testing.T) {
		adapter, err := NewConfigAdapter(GoConfigsSource)
		if err != nil {
			t.Fatalf("Failed to create adapter: %v", err)
		}
		
		os.Remove(filepath.Join(goConfigsDir, "prompts.txt"))
		
		prompts, err := adapter.LoadPrompts()
		if err != nil {
			t.Fatalf("Failed to load prompts with fallback: %v", err)
		}
		
		if prompts != pkgPromptsContent {
			t.Errorf("Expected fallback content %q, got %q", pkgPromptsContent, prompts)
		}
	})
	
	t.Run("Test UpdateGoConfigsFromPkg", func(t *testing.T) {
		adapter, err := NewConfigAdapter(PkgSource)
		if err != nil {
			t.Fatalf("Failed to create adapter: %v", err)
		}
		
		pkgPromptsWithSamSepiol := "# Test prompts from pkg\nYou are Sam Sepiol, a software engineer."
		os.WriteFile(filepath.Join(pkgPromptsDir, "prompts.txt"), []byte(pkgPromptsWithSamSepiol), 0644)
		
		err = adapter.UpdateGoConfigsFromPkg()
		if err != nil {
			t.Fatalf("Failed to update go_configs from pkg: %v", err)
		}
		
		updatedGoConfigsPrompts, err := os.ReadFile(filepath.Join(goConfigsDir, "prompts.txt"))
		if err != nil {
			t.Fatalf("Failed to read updated go_configs prompts.txt: %v", err)
		}
		
		if !strings.Contains(string(updatedGoConfigsPrompts), "samsepi0l") {
			t.Errorf("Expected 'samsepi0l' in updated go_configs prompts.txt, got %q", string(updatedGoConfigsPrompts))
		}
		
		updatedGoConfigsTools, err := os.ReadFile(filepath.Join(goConfigsDir, "tools.json"))
		if err != nil {
			t.Fatalf("Failed to read updated go_configs tools.json: %v", err)
		}
		
		var toolsConfig map[string]interface{}
		if err := json.Unmarshal(updatedGoConfigsTools, &toolsConfig); err != nil {
			t.Fatalf("Failed to parse updated go_configs tools.json: %v", err)
		}
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if !ok || len(toolsArray) == 0 {
			t.Fatalf("Invalid tools format or empty tools array in updated go_configs tools.json")
		}
		
		firstTool, ok := toolsArray[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid tool format in updated go_configs tools.json")
		}
		
		if firstTool["description"] != "A test tool from pkg" {
			t.Errorf("Expected tool description 'A test tool from pkg' in updated go_configs tools.json, got %v", firstTool["description"])
		}
	})
	
	t.Run("Test UpdatePkgFromGoConfigs", func(t *testing.T) {
		adapter, err := NewConfigAdapter(GoConfigsSource)
		if err != nil {
			t.Fatalf("Failed to create adapter: %v", err)
		}
		
		goConfigsPromptsWithSamSepiol := "# Test prompts from go_configs\nYou are Sam Sepiol, a software engineer."
		os.WriteFile(filepath.Join(goConfigsDir, "prompts.txt"), []byte(goConfigsPromptsWithSamSepiol), 0644)
		
		err = adapter.UpdatePkgFromGoConfigs()
		if err != nil {
			t.Fatalf("Failed to update pkg from go_configs: %v", err)
		}
		
		updatedPkgPrompts, err := os.ReadFile(filepath.Join(pkgPromptsDir, "prompts.txt"))
		if err != nil {
			t.Fatalf("Failed to read updated pkg prompts.txt: %v", err)
		}
		
		if !strings.Contains(string(updatedPkgPrompts), "samsepi0l") {
			t.Errorf("Expected 'samsepi0l' in updated pkg prompts.txt, got %q", string(updatedPkgPrompts))
		}
		
		updatedPkgTools, err := os.ReadFile(filepath.Join(pkgToolsDir, "tools.json"))
		if err != nil {
			t.Fatalf("Failed to read updated pkg tools.json: %v", err)
		}
		
		var toolsConfig map[string]interface{}
		if err := json.Unmarshal(updatedPkgTools, &toolsConfig); err != nil {
			t.Fatalf("Failed to parse updated pkg tools.json: %v", err)
		}
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if !ok || len(toolsArray) == 0 {
			t.Fatalf("Invalid tools format or empty tools array in updated pkg tools.json")
		}
		
		firstTool, ok := toolsArray[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Invalid tool format in updated pkg tools.json")
		}
		
		if firstTool["description"] != "A test tool" {
			t.Errorf("Expected tool description 'A test tool' in updated pkg tools.json, got %v", firstTool["description"])
		}
	})
}
