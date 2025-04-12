package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigInteroperability(t *testing.T) {
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
	
	goConfigsToolsContent := `{
		"tools": [
			{
				"name": "test_tool",
				"description": "A test tool",
				"parameters": {
					"param1": {
						"type": "string",
						"description": "An example parameter."
					}
				},
				"returns": {
					"type": "string",
					"description": "Operation result"
				}
			}
		]
	}`
	
	pkgToolsContent := `{
		"tools": [
			{
				"name": "test_tool",
				"description": "A test tool from pkg",
				"parameters": {
					"param1": {
						"type": "string",
						"description": "An example parameter."
					}
				},
				"returns": {
					"type": "string",
					"description": "Operation result"
				}
			}
		]
	}`
	
	os.WriteFile(filepath.Join(goConfigsDir, "prompts.txt"), []byte(goConfigsPromptsContent), 0644)
	os.WriteFile(filepath.Join(pkgPromptsDir, "prompts.txt"), []byte(pkgPromptsContent), 0644)
	os.WriteFile(filepath.Join(goConfigsDir, "tools.json"), []byte(goConfigsToolsContent), 0644)
	os.WriteFile(filepath.Join(pkgToolsDir, "tools.json"), []byte(pkgToolsContent), 0644)
	
	t.Run("Test loading tools from both locations", func(t *testing.T) {
		goConfigsData, err := os.ReadFile(filepath.Join(goConfigsDir, "tools.json"))
		if err != nil {
			t.Fatalf("Failed to read go_configs tools.json: %v", err)
		}
		
		var goConfigsTools map[string]interface{}
		if err := json.Unmarshal(goConfigsData, &goConfigsTools); err != nil {
			t.Fatalf("Failed to parse go_configs tools.json: %v", err)
		}
		
		pkgData, err := os.ReadFile(filepath.Join(pkgToolsDir, "tools.json"))
		if err != nil {
			t.Fatalf("Failed to read pkg tools.json: %v", err)
		}
		
		var pkgTools map[string]interface{}
		if err := json.Unmarshal(pkgData, &pkgTools); err != nil {
			t.Fatalf("Failed to parse pkg tools.json: %v", err)
		}
		
		goConfigsToolsArr, ok := goConfigsTools["tools"].([]interface{})
		if !ok || len(goConfigsToolsArr) == 0 {
			t.Fatal("Invalid tools array structure in go_configs tools.json")
		}
		
		pkgToolsArr, ok := pkgTools["tools"].([]interface{})
		if !ok || len(pkgToolsArr) == 0 {
			t.Fatal("Invalid tools array structure in pkg tools.json")
		}
		
		goConfigsTool := goConfigsToolsArr[0].(map[string]interface{})
		pkgTool := pkgToolsArr[0].(map[string]interface{})
		
		requiredFields := []string{"name", "description", "parameters", "returns"}
		for _, field := range requiredFields {
			if _, exists := goConfigsTool[field]; !exists {
				t.Errorf("Missing required field %s in go_configs tool", field)
			}
			if _, exists := pkgTool[field]; !exists {
				t.Errorf("Missing required field %s in pkg tool", field)
			}
		}
	})
	
	t.Run("Test loading prompts from both locations", func(t *testing.T) {
		goConfigsPrompts, err := os.ReadFile(filepath.Join(goConfigsDir, "prompts.txt"))
		if err != nil {
			t.Fatalf("Failed to read go_configs prompts.txt: %v", err)
		}
		
		pkgPrompts, err := os.ReadFile(filepath.Join(pkgPromptsDir, "prompts.txt"))
		if err != nil {
			t.Fatalf("Failed to read pkg prompts.txt: %v", err)
		}
		
		if !strings.Contains(string(goConfigsPrompts), "samsepi0l") {
			t.Error("go_configs prompts.txt does not contain required agent name")
		}
		if !strings.Contains(string(pkgPrompts), "samsepi0l") {
			t.Error("pkg prompts.txt does not contain required agent name")
		}
	})
}
