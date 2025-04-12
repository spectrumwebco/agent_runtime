package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func SyncConfigs(direction string) error {
	var adapter *ConfigAdapter
	var err error
	
	switch direction {
	case "pkg-to-go":
		adapter, err = NewConfigAdapter(PkgSource)
		if err != nil {
			return fmt.Errorf("failed to create config adapter: %w", err)
		}
		
		fmt.Println("Synchronizing configurations from pkg to go_configs...")
		err = adapter.UpdateGoConfigsFromPkg()
		if err != nil {
			return fmt.Errorf("failed to update go_configs from pkg: %w", err)
		}
		
	case "go-to-pkg":
		adapter, err = NewConfigAdapter(GoConfigsSource)
		if err != nil {
			return fmt.Errorf("failed to create config adapter: %w", err)
		}
		
		fmt.Println("Synchronizing configurations from go_configs to pkg...")
		err = adapter.UpdatePkgFromGoConfigs()
		if err != nil {
			return fmt.Errorf("failed to update pkg from go_configs: %w", err)
		}
		
	default:
		return fmt.Errorf("invalid direction: %s. Must be 'pkg-to-go' or 'go-to-pkg'", direction)
	}
	
	fmt.Println("Configuration synchronization completed successfully.")
	return nil
}

func EnsureConfigDirectories() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return fmt.Errorf("failed to find repository root: %w", err)
	}
	
	goConfigsDir := filepath.Join(repoRoot, "go_configs")
	if err := os.MkdirAll(goConfigsDir, 0755); err != nil {
		return fmt.Errorf("failed to create go_configs directory: %w", err)
	}
	
	pkgPromptsDir := filepath.Join(repoRoot, "pkg", "prompts")
	if err := os.MkdirAll(pkgPromptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create pkg/prompts directory: %w", err)
	}
	
	pkgToolsDir := filepath.Join(repoRoot, "pkg", "tools")
	if err := os.MkdirAll(pkgToolsDir, 0755); err != nil {
		return fmt.Errorf("failed to create pkg/tools directory: %w", err)
	}
	
	return nil
}
