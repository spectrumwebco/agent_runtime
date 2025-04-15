package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent   AgentConfig   `yaml:"agent"`
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	LLM     LLMConfig     `yaml:"llm"`
	MCP     MCPConfig     `yaml:"mcp"`
	Python  PythonConfig  `yaml:"python"`
	CPP     CPPConfig     `yaml:"cpp"`
	Runtime RuntimeConfig `yaml:"runtime"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	GRPCPort int    `yaml:"grpcPort"`
	GRPCHost string `yaml:"grpcHost"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	File   string `yaml:"file"`
}

type LLMConfig struct {
	Provider string            `yaml:"provider"`
	Model    string            `yaml:"model"`
	Options  map[string]string `yaml:"options"`
}

type MCPConfig struct {
	Servers []MCPServerConfig `yaml:"servers"`
}

type MCPServerConfig struct {
	Name    string            `yaml:"name"`
	Enabled bool              `yaml:"enabled"`
	Options map[string]string `yaml:"options"`
}

type PythonConfig struct {
	Enabled     bool     `yaml:"enabled"`
	Interpreter string   `yaml:"interpreter"`
	Paths       []string `yaml:"paths"`
}

type CPPConfig struct {
	Enabled     bool              `yaml:"enabled"`
	Compiler    string            `yaml:"compiler"`
	Flags       []string          `yaml:"flags"`
	Libraries   []string          `yaml:"libraries"`
	Options     map[string]string `yaml:"options"`
}

type AgentConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type RuntimeConfig struct {
	Sandbox      bool              `yaml:"sandbox"`
	AllowedPaths []string          `yaml:"allowedPaths"`
	Environment  map[string]string `yaml:"environment"`
}

func DefaultConfig() *Config {
	return &Config{
		Agent: AgentConfig{
			Name:        "Sam Sepiol",
			Description: "An autonomous software engineering agent",
			Version:     "0.1.0",
		},
		Server: ServerConfig{
			Port:     8080,
			Host:     "0.0.0.0",
			GRPCPort: 50051,
			GRPCHost: "0.0.0.0",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		LLM: LLMConfig{
			Provider: "openai",
			Model:    "gpt-4",
		},
		MCP: MCPConfig{
			Servers: []MCPServerConfig{
				{
					Name:    "filesystem",
					Enabled: true,
				},
				{
					Name:    "tools",
					Enabled: true,
				},
			},
		},
		Python: PythonConfig{
			Enabled:     true,
			Interpreter: "python3",
		},
		CPP: CPPConfig{
			Enabled:  true,
			Compiler: "g++",
			Flags:    []string{"-std=c++17", "-O2"},
			Libraries: []string{
				"boost",
				"eigen",
			},
		},
		Runtime: RuntimeConfig{
			Sandbox: true,
		},
	}
}

func Load(path string) (*Config, error) {
	config := DefaultConfig()

	if path == "" {
		return config, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", path)
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return config, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(filepath.Clean(path), data, 0600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}
