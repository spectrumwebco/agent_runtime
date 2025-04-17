package config

import "fmt"

type AgentConfig struct {
	Name                 string `mapstructure:"name" yaml:"name"`
	MaxRequeries         int    `mapstructure:"max_requeries" yaml:"max_requeries"`
	AlwaysRequireZeroExit bool   `mapstructure:"always_require_zero_exit" yaml:"always_require_zero_exit"`
	ExecutionTimeout     int    `mapstructure:"execution_timeout" yaml:"execution_timeout"`
	ConfigSource         string `mapstructure:"config_source" yaml:"config_source"`
}

type ToolsConfiguration struct {
	ExecutionTimeout      int `mapstructure:"execution_timeout" yaml:"execution_timeout"`
	MaxOutputSize         int `mapstructure:"max_output_size" yaml:"max_output_size"`
	TotalExecutionTimeout int `mapstructure:"total_execution_timeout" yaml:"total_execution_timeout"`
	MaxConsecutiveTimeouts int `mapstructure:"max_consecutive_timeouts" yaml:"max_consecutive_timeouts"`
}

type Config struct {
	ServerPort string         `mapstructure:"server_port" yaml:"server_port"`
	LogLevel   string         `mapstructure:"log_level" yaml:"log_level"`
	Agent      AgentConfig    `mapstructure:"agent" yaml:"agent"`
	Tools      ToolsConfiguration    `mapstructure:"tools" yaml:"tools"`
	RocketMQ   RocketMQConfig `mapstructure:"rocketmq" yaml:"rocketmq"`
	Supabase   SupabaseConfig `mapstructure:"supabase" yaml:"supabase"`
	Dragonfly  DragonflyConfig `mapstructure:"dragonfly" yaml:"dragonfly"`
	AI         AIConfig       `mapstructure:"ai" yaml:"ai"`
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

type DragonflyConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
}

type AIConfig struct {
	DefaultProvider string            `mapstructure:"default_provider" yaml:"default_provider"`
	Providers       []string          `mapstructure:"providers" yaml:"providers"`
	APIKeys         map[string]string `mapstructure:"api_keys" yaml:"api_keys"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerPort: "8080",
		LogLevel:   "info",
		Agent: AgentConfig{
			Name:                 "samsepi0l",
			MaxRequeries:         3,
			AlwaysRequireZeroExit: true,
			ExecutionTimeout:     60,
			ConfigSource:         "pkg",
		},
		Tools: ToolsConfiguration{
			ExecutionTimeout:      60,
			MaxOutputSize:         10000,
			TotalExecutionTimeout: 1800,
			MaxConsecutiveTimeouts: 3,
		},
		RocketMQ: RocketMQConfig{
			NameServerAddrs:      []string{"127.0.0.1:9876"},
			StateTopic:           "agent_state_updates",
			ProducerGroupName:    "agent_runtime_producer_group",
			ConsumerGroupName:    "agent_runtime_consumer_group",
			ConsumerSubscription: "*",
			ProducerRetry:        2,
		},
		Supabase: SupabaseConfig{
			MainURL:     "https://supabase.example.com",
			ReadonlyURL: "https://readonly.supabase.example.com",
			RollbackURL: "https://rollback.supabase.example.com",
			APIKey:      "",
			AuthToken:   "",
		},
		Dragonfly: DragonflyConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
		},
		AI: AIConfig{
			DefaultProvider: "klusterai",
			Providers:       []string{"klusterai", "openrouter", "bedrock", "gemini"},
			APIKeys:         make(map[string]string),
		},
	}
}

func Load() (*Config, error) {
	fmt.Println("Placeholder: Loading configuration...")
	return &Config{
		ServerPort: "8080", // Default example
		LogLevel:   "info",  // Default example
	}, nil
}
