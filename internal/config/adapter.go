package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type ConfigSource string

const (
	GoConfigsSource ConfigSource = "go_configs"
	
	PkgSource ConfigSource = "pkg"
)

type ConfigType string

const (
	PromptsConfig ConfigType = "prompts.txt"
	
	ToolsConfig ConfigType = "tools.json"
)

type ConfigAdapter struct {
	RepoRoot string
	
	PreferredSource ConfigSource
}

func NewConfigAdapter(preferredSource ConfigSource) (*ConfigAdapter, error) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return nil, err
	}
	
	return &ConfigAdapter{
		RepoRoot:        repoRoot,
		PreferredSource: preferredSource,
	}, nil
}

func (a *ConfigAdapter) LoadPrompts() (string, error) {
	data, err := a.LoadConfig(PromptsConfig)
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

func (a *ConfigAdapter) LoadTools() (map[string]interface{}, error) {
	data, err := a.LoadConfig(ToolsConfig)
	if err != nil {
		return nil, err
	}
	
	var toolsConfig map[string]interface{}
	if err := json.Unmarshal(data, &toolsConfig); err != nil {
		return nil, fmt.Errorf("failed to parse tools.json: %v", err)
	}
	
	return toolsConfig, nil
}

func (a *ConfigAdapter) LoadConfig(configType ConfigType) ([]byte, error) {
	var primaryPath, fallbackPath string
	
	switch a.PreferredSource {
	case GoConfigsSource:
		primaryPath = filepath.Join(a.RepoRoot, "go_configs", string(configType))
		fallbackPath = filepath.Join(a.RepoRoot, "pkg", getConfigSubdir(configType), string(configType))
	case PkgSource:
		primaryPath = filepath.Join(a.RepoRoot, "pkg", getConfigSubdir(configType), string(configType))
		fallbackPath = filepath.Join(a.RepoRoot, "go_configs", string(configType))
	default:
		return nil, fmt.Errorf("invalid preferred source: %s", a.PreferredSource)
	}
	
	data, err := ioutil.ReadFile(primaryPath)
	if err == nil {
		return data, nil
	}
	
	data, err = ioutil.ReadFile(fallbackPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from both sources: %v", err)
	}
	
	return data, nil
}

func getConfigSubdir(configType ConfigType) string {
	switch configType {
	case PromptsConfig:
		return "prompts"
	case ToolsConfig:
		return "tools"
	default:
		return ""
	}
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

func (a *ConfigAdapter) UpdateGoConfigsFromPkg() error {
	pkgPromptsPath := filepath.Join(a.RepoRoot, "pkg", "prompts", "prompts.txt")
	goConfigsPromptsPath := filepath.Join(a.RepoRoot, "go_configs", "prompts.txt")
	
	pkgPromptsData, err := ioutil.ReadFile(pkgPromptsPath)
	if err != nil {
		return fmt.Errorf("failed to read pkg prompts.txt: %v", err)
	}
	
	promptsContent := string(pkgPromptsData)
	if !strings.Contains(promptsContent, "samsepi0l") {
		promptsContent = strings.ReplaceAll(promptsContent, "Sam Sepiol", "samsepi0l")
	}
	
	if err := ioutil.WriteFile(goConfigsPromptsPath, []byte(promptsContent), 0644); err != nil {
		return fmt.Errorf("failed to write go_configs prompts.txt: %v", err)
	}
	
	pkgToolsPath := filepath.Join(a.RepoRoot, "pkg", "tools", "tools.json")
	goConfigsToolsPath := filepath.Join(a.RepoRoot, "go_configs", "tools.json")
	
	pkgToolsData, err := ioutil.ReadFile(pkgToolsPath)
	if err != nil {
		return fmt.Errorf("failed to read pkg tools.json: %v", err)
	}
	
	if err := ioutil.WriteFile(goConfigsToolsPath, pkgToolsData, 0644); err != nil {
		return fmt.Errorf("failed to write go_configs tools.json: %v", err)
	}
	
	return nil
}

func (a *ConfigAdapter) UpdatePkgFromGoConfigs() error {
	goConfigsPromptsPath := filepath.Join(a.RepoRoot, "go_configs", "prompts.txt")
	pkgPromptsPath := filepath.Join(a.RepoRoot, "pkg", "prompts", "prompts.txt")
	
	goConfigsPromptsData, err := ioutil.ReadFile(goConfigsPromptsPath)
	if err != nil {
		return fmt.Errorf("failed to read go_configs prompts.txt: %v", err)
	}
	
	promptsContent := string(goConfigsPromptsData)
	if !strings.Contains(promptsContent, "samsepi0l") {
		promptsContent = strings.ReplaceAll(promptsContent, "Sam Sepiol", "samsepi0l")
	}
	
	if err := ioutil.WriteFile(pkgPromptsPath, []byte(promptsContent), 0644); err != nil {
		return fmt.Errorf("failed to write pkg prompts.txt: %v", err)
	}
	
	goConfigsToolsPath := filepath.Join(a.RepoRoot, "go_configs", "tools.json")
	pkgToolsPath := filepath.Join(a.RepoRoot, "pkg", "tools", "tools.json")
	
	goConfigsToolsData, err := ioutil.ReadFile(goConfigsToolsPath)
	if err != nil {
		return fmt.Errorf("failed to read go_configs tools.json: %v", err)
	}
	
	if err := ioutil.WriteFile(pkgToolsPath, goConfigsToolsData, 0644); err != nil {
		return fmt.Errorf("failed to write pkg tools.json: %v", err)
	}
	
	return nil
}

func (a *ConfigAdapter) GetAgentName() (string, error) {
	prompts, err := a.LoadPrompts()
	if err != nil {
		return "samsepi0l", nil // Default agent name
	}
	
	if strings.Contains(prompts, "samsepi0l") {
		return "samsepi0l", nil
	}
	
	if strings.Contains(prompts, "Sam Sepiol") {
		return "Sam Sepiol", nil
	}
	
	if strings.Contains(prompts, "Kled") {
		return "Kled", nil
	}
	
	return "samsepi0l", nil
}
