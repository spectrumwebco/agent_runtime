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
	
	v.SetDefault("rocketmq.name_server_addrs", []string{"127.0.0.1:9876"}) // Default local nameserver
	v.SetDefault("rocketmq.state_topic", "agent_state_updates")
	v.SetDefault("rocketmq.producer_group_name", "agent_runtime_producer_group")
	v.SetDefault("rocketmq.consumer_group_name", "agent_runtime_consumer_group")
	v.SetDefault("rocketmq.consumer_subscription", "*") // Consume all tags by default
	v.SetDefault("rocketmq.producer_retry", 2)          // Default retry count for producer
	v.SetDefault("supabase.main_url", "https://supabase.example.com")
	v.SetDefault("supabase.readonly_url", "https://readonly.supabase.example.com")
	v.SetDefault("supabase.rollback_url", "https://rollback.supabase.example.com")
	v.SetDefault("supabase.api_key", "")
	v.SetDefault("supabase.auth_token", "")




	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist, using defaults\n", configPath)
		var defaultConfig Config
		if err := v.Unmarshal(&defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal default config: %w", err)
		}
		return &defaultConfig, nil
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
