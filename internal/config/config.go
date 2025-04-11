package config

import "fmt"

type Config struct {
	ServerPort string `yaml:"server_port"`
	LogLevel   string `yaml:"log_level"`
}

func Load() (*Config, error) {
	fmt.Println("Placeholder: Loading configuration...")
	return &Config{
		ServerPort: "8080", // Default example
		LogLevel:   "info",  // Default example
	}, nil
}
