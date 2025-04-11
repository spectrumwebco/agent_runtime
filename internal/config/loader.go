package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	
	v.SetDefault("agent.name", "default-go-agent")
	v.SetDefault("agent.max_requeries", 3)
	v.SetDefault("agent.always_require_zero_exit", true)
	v.SetDefault("tools.execution_timeout", 60) // 60 seconds
	v.SetDefault("tools.max_output_size", 10000) // 10000 characters
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist, using defaults\n", configPath)
		return &Config{
			Agent: AgentConfig{
				Name:                 "default-go-agent",
				MaxRequeries:         3,
				AlwaysRequireZeroExit: true,
			},
			Tools: ToolsConfig{
				ExecutionTimeout: 60,
				MaxOutputSize:    10000,
			},
		}, nil
	}
	
	dir, file := filepath.Split(configPath)
	ext := filepath.Ext(file)
	name := file[:len(file)-len(ext)]
	
	v.SetConfigName(name)
	v.SetConfigType(ext[1:]) // Remove the dot from extension
	v.AddConfigPath(dir)
	
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}
