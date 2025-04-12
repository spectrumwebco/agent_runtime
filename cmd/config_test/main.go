package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/pkg/modules"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

func main() {
	fmt.Println("Testing configuration interchangeability...")
	
	fmt.Println("\nTesting loading from go_configs:")
	adapter, err := config.NewConfigAdapter(config.GoConfigsSource)
	if err != nil {
		fmt.Printf("Error creating config adapter: %v\n", err)
		os.Exit(1)
	}
	
	prompts, err := adapter.LoadPrompts()
	if err != nil {
		fmt.Printf("Error loading prompts from go_configs: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded prompts from go_configs (%d bytes)\n", len(prompts))
	}
	
	toolsConfig, err := adapter.LoadTools()
	if err != nil {
		fmt.Printf("Error loading tools from go_configs: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded tools from go_configs\n")
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if ok {
			fmt.Printf("Found %d tools in go_configs/tools.json\n", len(toolsArray))
		}
	}
	
	fmt.Println("\nTesting loading from pkg:")
	adapter, err = config.NewConfigAdapter(config.PkgSource)
	if err != nil {
		fmt.Printf("Error creating config adapter: %v\n", err)
		os.Exit(1)
	}
	
	prompts, err = adapter.LoadPrompts()
	if err != nil {
		fmt.Printf("Error loading prompts from pkg: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded prompts from pkg (%d bytes)\n", len(prompts))
	}
	
	toolsConfig, err = adapter.LoadTools()
	if err != nil {
		fmt.Printf("Error loading tools from pkg: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded tools from pkg\n")
		
		toolsArray, ok := toolsConfig["tools"].([]interface{})
		if ok {
			fmt.Printf("Found %d tools in pkg/tools/tools.json\n", len(toolsArray))
		}
	}
	
	fmt.Println("\nTesting configuration synchronization:")
	err = adapter.UpdateGoConfigsFromPkg()
	if err != nil {
		fmt.Printf("Error synchronizing from pkg to go_configs: %v\n", err)
	} else {
		fmt.Println("Successfully synchronized configurations from pkg to go_configs")
	}
	
	fmt.Println("\nTesting module registry with tools:")
	registry := modules.NewRegistry()
	
	module := modules.NewBaseModule("test_module", "Test module for configuration interchangeability")
	registry.Register(module)
	
	fmt.Printf("Registered %d modules\n", len(registry.List()))
	
	adapter, _ = config.NewConfigAdapter(config.GoConfigsSource)
	toolsConfig, _ = adapter.LoadTools()
	
	toolRegistry := &tools.ToolRegistry{}
	err = tools.LoadToolsFromConfig(toolRegistry)
	if err != nil {
		fmt.Printf("Error loading tools from config: %v\n", err)
	} else {
		fmt.Printf("Successfully loaded tools into registry\n")
		fmt.Printf("Tool registry has %d commands\n", len(toolRegistry.ListCommands()))
	}
	
	fmt.Println("\nConfiguration interchangeability test completed")
}
