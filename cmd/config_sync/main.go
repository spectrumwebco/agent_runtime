package main

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/spectrumwebco/agent_runtime/internal/config"
)

func main() {
	direction := flag.String("direction", "pkg-to-go", "Direction of synchronization: pkg-to-go or go-to-pkg")
	flag.Parse()
	
	var adapter *config.ConfigAdapter
	var err error
	
	switch *direction {
	case "pkg-to-go":
		adapter, err = config.NewConfigAdapter(config.PkgSource)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config adapter: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Synchronizing configurations from pkg to go_configs...")
		err = adapter.UpdateGoConfigsFromPkg()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update go_configs from pkg: %v\n", err)
			os.Exit(1)
		}
		
	case "go-to-pkg":
		adapter, err = config.NewConfigAdapter(config.GoConfigsSource)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config adapter: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("Synchronizing configurations from go_configs to pkg...")
		err = adapter.UpdatePkgFromGoConfigs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update pkg from go_configs: %v\n", err)
			os.Exit(1)
		}
		
	default:
		fmt.Fprintf(os.Stderr, "Invalid direction: %s. Must be 'pkg-to-go' or 'go-to-pkg'\n", *direction)
		os.Exit(1)
	}
	
	fmt.Println("Configuration synchronization completed successfully.")
}
