package main

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/spectrumwebco/agent_runtime/internal/config"
)

func main() {
	configType := flag.String("type", "prompts", "Type of configuration to load (prompts or tools)")
	location := flag.String("location", "pkg", "Location to load from (go_configs or pkg)")
	outputFile := flag.String("output", "", "Output file path (optional)")
	
	flag.Parse()
	
	var configTypeEnum config.ConfigType
	switch *configType {
	case "prompts":
		configTypeEnum = config.PromptsConfig
	case "tools":
		configTypeEnum = config.ToolsConfig
	default:
		fmt.Fprintf(os.Stderr, "Invalid config type: %s. Must be 'prompts' or 'tools'\n", *configType)
		os.Exit(1)
	}
	
	var locationEnum config.ConfigLocation
	switch *location {
	case "go_configs":
		locationEnum = config.GoConfigsLocation
	case "pkg":
		locationEnum = config.PkgLocation
	default:
		fmt.Fprintf(os.Stderr, "Invalid location: %s. Must be 'go_configs' or 'pkg'\n", *location)
		os.Exit(1)
	}
	
	data, err := config.LoadConfig(configTypeEnum, locationEnum, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	
	if *outputFile == "" {
		fmt.Println(string(data))
	} else {
		err := os.WriteFile(*outputFile, data, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Configuration written to %s\n", *outputFile)
	}
}
