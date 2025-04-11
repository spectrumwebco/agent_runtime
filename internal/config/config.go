package config

import "fmt"

type AgentConfig struct {
	Name                 string `mapstructure:"name" yaml:"name"`
	MaxRequeries         int    `mapstructure:"max_requeries" yaml:"max_requeries"`
	AlwaysRequireZeroExit bool   `mapstructure:"always_require_zero_exit" yaml:"always_require_zero_exit"`
}

type ToolsConfig struct {
	ExecutionTimeout int `mapstructure:"execution_timeout" yaml:"execution_timeout"`
	MaxOutputSize    int `mapstructure:"max_output_size" yaml:"max_output_size"`
}

type Config struct {
	ServerPort string         `mapstructure:"server_port" yaml:"server_port"`
	LogLevel   string         `mapstructure:"log_level" yaml:"log_level"`
	Agent      AgentConfig    `mapstructure:"agent" yaml:"agent"`
	Tools      ToolsConfig    `mapstructure:"tools" yaml:"tools"`
	RocketMQ   RocketMQConfig `mapstructure:"rocketmq" yaml:"rocketmq"`
	Supabase   SupabaseConfig `mapstructure:"supabase" yaml:"supabase"`
}
type RocketMQConfig struct {
	NameServerAddrs        []string `mapstructure:"name_server_addrs" yaml:"name_server_addrs"`
	StateTopic             string   `mapstructure:"state_topic" yaml:"state_topic"`
	ProducerGroupName      string   `mapstructure:"producer_group_name" yaml:"producer_group_name"`
	ConsumerGroupName      string   `mapstructure:"consumer_group_name" yaml:"consumer_group_name"`
	ConsumerSubscription   string   `mapstructure:"consumer_subscription" yaml:"consumer_subscription"` // e.g., "*" or "TagA"
	ProducerRetry          int      `mapstructure:"producer_retry" yaml:"producer_retry"`
}

type SupabaseConfig struct {
	MainURL       string `mapstructure:"main_url" yaml:"main_url"`
	ReadonlyURL   string `mapstructure:"readonly_url" yaml:"readonly_url"`
	RollbackURL   string `mapstructure:"rollback_url" yaml:"rollback_url"`
	APIKey        string `mapstructure:"api_key" yaml:"api_key"`
	AuthToken     string `mapstructure:"auth_token" yaml:"auth_token"`
}



func Load() (*Config, error) {
	fmt.Println("Placeholder: Loading configuration...")
	return &Config{
		ServerPort: "8080", // Default example
		LogLevel:   "info",  // Default example
	}, nil
}
