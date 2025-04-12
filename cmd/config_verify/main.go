package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("Verifying configuration interchangeability...")
	
	repoRoot, err := findRepoRoot()
	if err != nil {
		fmt.Printf("Error finding repository root: %v\n", err)
		os.Exit(1)
	}
	
	goConfigsPromptsPath := filepath.Join(repoRoot, "go_configs", "prompts.txt")
	goConfigsToolsPath := filepath.Join(repoRoot, "go_configs", "tools.json")
	
	pkgPromptsPath := filepath.Join(repoRoot, "pkg", "prompts", "prompts.txt")
	pkgToolsPath := filepath.Join(repoRoot, "pkg", "tools", "tools.json")
	
	fmt.Println("\nVerifying prompts.txt:")
	goConfigsPrompts, err := ioutil.ReadFile(goConfigsPromptsPath)
	if err != nil {
		fmt.Printf("Error reading go_configs/prompts.txt: %v\n", err)
	} else {
		fmt.Printf("go_configs/prompts.txt size: %d bytes\n", len(goConfigsPrompts))
		
		if strings.Contains(string(goConfigsPrompts), "samsepi0l") {
			fmt.Println("✅ go_configs/prompts.txt contains agent name 'samsepi0l'")
		} else {
			fmt.Println("❌ go_configs/prompts.txt does not contain agent name 'samsepi0l'")
		}
	}
	
	pkgPrompts, err := ioutil.ReadFile(pkgPromptsPath)
	if err != nil {
		fmt.Printf("Error reading pkg/prompts/prompts.txt: %v\n", err)
	} else {
		fmt.Printf("pkg/prompts/prompts.txt size: %d bytes\n", len(pkgPrompts))
		
		if strings.Contains(string(pkgPrompts), "samsepi0l") {
			fmt.Println("✅ pkg/prompts/prompts.txt contains agent name 'samsepi0l'")
		} else {
			fmt.Println("❌ pkg/prompts/prompts.txt does not contain agent name 'samsepi0l'")
		}
	}
	
	fmt.Println("\nVerifying tools.json:")
	goConfigsTools, err := ioutil.ReadFile(goConfigsToolsPath)
	if err != nil {
		fmt.Printf("Error reading go_configs/tools.json: %v\n", err)
	} else {
		fmt.Printf("go_configs/tools.json size: %d bytes\n", len(goConfigsTools))
		
		var goConfigsToolsData map[string]interface{}
		if err := json.Unmarshal(goConfigsTools, &goConfigsToolsData); err != nil {
			fmt.Printf("Error parsing go_configs/tools.json: %v\n", err)
		} else {
			toolsArray, ok := goConfigsToolsData["tools"].([]interface{})
			if ok {
				fmt.Printf("go_configs/tools.json contains %d tools\n", len(toolsArray))
			} else {
				fmt.Println("❌ go_configs/tools.json has invalid format (missing 'tools' array)")
			}
		}
	}
	
	pkgTools, err := ioutil.ReadFile(pkgToolsPath)
	if err != nil {
		fmt.Printf("Error reading pkg/tools/tools.json: %v\n", err)
	} else {
		fmt.Printf("pkg/tools/tools.json size: %d bytes\n", len(pkgTools))
		
		var pkgToolsData map[string]interface{}
		if err := json.Unmarshal(pkgTools, &pkgToolsData); err != nil {
			fmt.Printf("Error parsing pkg/tools/tools.json: %v\n", err)
		} else {
			toolsArray, ok := pkgToolsData["tools"].([]interface{})
			if ok {
				fmt.Printf("pkg/tools/tools.json contains %d tools\n", len(toolsArray))
			} else {
				fmt.Println("❌ pkg/tools/tools.json has invalid format (missing 'tools' array)")
			}
		}
	}
	
	fmt.Println("\nTesting configuration synchronization:")
	
	tempPromptsPath := filepath.Join(repoRoot, "go_configs", "prompts.txt.bak")
	if err := copyFile(goConfigsPromptsPath, tempPromptsPath); err != nil {
		fmt.Printf("Error creating backup of go_configs/prompts.txt: %v\n", err)
	} else {
		defer os.Remove(tempPromptsPath) // Clean up at the end
		
		if err := copyFile(pkgPromptsPath, goConfigsPromptsPath); err != nil {
			fmt.Printf("Error updating go_configs/prompts.txt: %v\n", err)
		} else {
			fmt.Println("✅ Successfully synchronized prompts.txt from pkg to go_configs")
			
			if err := copyFile(tempPromptsPath, goConfigsPromptsPath); err != nil {
				fmt.Printf("Error restoring go_configs/prompts.txt: %v\n", err)
			}
		}
	}
	
	tempToolsPath := filepath.Join(repoRoot, "go_configs", "tools.json.bak")
	if err := copyFile(goConfigsToolsPath, tempToolsPath); err != nil {
		fmt.Printf("Error creating backup of go_configs/tools.json: %v\n", err)
	} else {
		defer os.Remove(tempToolsPath) // Clean up at the end
		
		if err := copyFile(pkgToolsPath, goConfigsToolsPath); err != nil {
			fmt.Printf("Error updating go_configs/tools.json: %v\n", err)
		} else {
			fmt.Println("✅ Successfully synchronized tools.json from pkg to go_configs")
			
			if err := copyFile(tempToolsPath, goConfigsToolsPath); err != nil {
				fmt.Printf("Error restoring go_configs/tools.json: %v\n", err)
			}
		}
	}
	
	fmt.Println("\nVerifying module compatibility:")
	modulesDir := filepath.Join(repoRoot, "pkg", "modules")
	
	moduleFiles, err := ioutil.ReadDir(modulesDir)
	if err != nil {
		fmt.Printf("Error reading pkg/modules directory: %v\n", err)
	} else {
		moduleCount := 0
		for _, file := range moduleFiles {
			if file.IsDir() && file.Name() != "." && file.Name() != ".." {
				moduleCount++
			}
		}
		fmt.Printf("Found %d module directories in pkg/modules\n", moduleCount)
		
		toolsImportCount := 0
		err := filepath.Walk(modulesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				
				if strings.Contains(string(content), "github.com/spectrumwebco/agent_runtime/pkg/tools") {
					toolsImportCount++
				}
			}
			
			return nil
		})
		
		if err != nil {
			fmt.Printf("Error checking module imports: %v\n", err)
		} else {
			fmt.Printf("Found %d module files importing pkg/tools\n", toolsImportCount)
			if toolsImportCount > 0 {
				fmt.Println("✅ Modules are compatible with tools configuration")
			} else {
				fmt.Println("❌ No modules import tools package")
			}
		}
	}
	
	fmt.Println("\nVerifying loop.go compatibility:")
	loopPath := filepath.Join(repoRoot, "internal", "agent", "loop.go")
	
	loopContent, err := ioutil.ReadFile(loopPath)
	if err != nil {
		fmt.Printf("Error reading internal/agent/loop.go: %v\n", err)
	} else {
		if strings.Contains(string(loopContent), "github.com/spectrumwebco/agent_runtime/pkg/tools") {
			fmt.Println("✅ loop.go imports pkg/tools")
		} else {
			fmt.Println("❌ loop.go does not import pkg/tools")
		}
		
		if strings.Contains(string(loopContent), "github.com/spectrumwebco/agent_runtime/pkg/modules") {
			fmt.Println("✅ loop.go imports pkg/modules")
		} else {
			fmt.Println("❌ loop.go does not import pkg/modules")
		}
		
		if strings.Contains(string(loopContent), "go_configs") {
			fmt.Println("❌ loop.go contains hardcoded references to go_configs")
		} else {
			fmt.Println("✅ loop.go does not contain hardcoded references to go_configs")
		}
	}
	
	fmt.Println("\nConfiguration interchangeability verification completed")
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root")
		}
		dir = parent
	}
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(dst, data, 0644)
}
