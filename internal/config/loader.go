package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	
	defaultConfig := NewDefaultConfig()
	
	v.SetDefault("agent.name", defaultConfig.Agent.Name)
	v.SetDefault("agent.max_requeries", defaultConfig.Agent.MaxRequeries)
	v.SetDefault("agent.always_require_zero_exit", defaultConfig.Agent.AlwaysRequireZeroExit)
	v.SetDefault("agent.execution_timeout", defaultConfig.Agent.ExecutionTimeout)
	v.SetDefault("agent.config_source", defaultConfig.Agent.ConfigSource)
	
	v.SetDefault("tools.execution_timeout", defaultConfig.Tools.ExecutionTimeout)
	v.SetDefault("tools.max_output_size", defaultConfig.Tools.MaxOutputSize)
	v.SetDefault("tools.total_execution_timeout", defaultConfig.Tools.TotalExecutionTimeout)
	v.SetDefault("tools.max_consecutive_timeouts", defaultConfig.Tools.MaxConsecutiveTimeouts)
	
	v.SetDefault("rocketmq.name_server_addrs", defaultConfig.RocketMQ.NameServerAddrs)
	v.SetDefault("rocketmq.state_topic", defaultConfig.RocketMQ.StateTopic)
	v.SetDefault("rocketmq.producer_group_name", defaultConfig.RocketMQ.ProducerGroupName)
	v.SetDefault("rocketmq.consumer_group_name", defaultConfig.RocketMQ.ConsumerGroupName)
	v.SetDefault("rocketmq.consumer_subscription", defaultConfig.RocketMQ.ConsumerSubscription)
	v.SetDefault("rocketmq.producer_retry", defaultConfig.RocketMQ.ProducerRetry)
	
	v.SetDefault("supabase.main_url", defaultConfig.Supabase.MainURL)
	v.SetDefault("supabase.readonly_url", defaultConfig.Supabase.ReadonlyURL)
	v.SetDefault("supabase.rollback_url", defaultConfig.Supabase.RollbackURL)
	v.SetDefault("supabase.api_key", defaultConfig.Supabase.APIKey)
	v.SetDefault("supabase.auth_token", defaultConfig.Supabase.AuthToken)
	
	v.SetDefault("dragonfly.host", defaultConfig.Dragonfly.Host)
	v.SetDefault("dragonfly.port", defaultConfig.Dragonfly.Port)
	v.SetDefault("dragonfly.password", defaultConfig.Dragonfly.Password)
	
	v.SetDefault("ai.default_provider", defaultConfig.AI.DefaultProvider)
	v.SetDefault("ai.providers", defaultConfig.AI.Providers)
	v.SetDefault("ai.api_keys", defaultConfig.AI.APIKeys)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist, using defaults\n", configPath)
		return defaultConfig, nil
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
	
	adapter, err := NewConfigAdapter(ConfigSource(config.Agent.ConfigSource))
	if err != nil {
		fmt.Printf("Warning: Failed to create config adapter: %v. Using default configurations.\n", err)
		return &config, nil
	}
	
	_, err = adapter.LoadPrompts()
	if err != nil {
		fmt.Printf("Warning: Failed to load prompts: %v. Using default prompts.\n", err)
	}
	
	_, err = adapter.LoadTools()
	if err != nil {
		fmt.Printf("Warning: Failed to load tools: %v. Using default tools.\n", err)
	}
	
	return &config, nil
}
