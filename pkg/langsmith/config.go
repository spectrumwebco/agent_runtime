package langsmith

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type LangSmithConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	APIKey string `json:"api_key" yaml:"api_key"`

	APIUrl string `json:"api_url" yaml:"api_url"`

	ProjectName string `json:"project_name" yaml:"project_name"`

	SelfHosted bool `json:"self_hosted" yaml:"self_hosted"`

	SelfHostedConfig *SelfHostedConfig `json:"self_hosted_config,omitempty" yaml:"self_hosted_config,omitempty"`
}

type SelfHostedConfig struct {
	URL string `json:"url" yaml:"url"`

	AdminUsername string `json:"admin_username" yaml:"admin_username"`

	AdminPassword string `json:"admin_password" yaml:"admin_password"`

	LicenseKey string `json:"license_key" yaml:"license_key"`

	DatabaseURL string `json:"database_url" yaml:"database_url"`

	RedisURL string `json:"redis_url" yaml:"redis_url"`

	StorageType string `json:"storage_type" yaml:"storage_type"`

	StorageConfig map[string]interface{} `json:"storage_config,omitempty" yaml:"storage_config,omitempty"`
}

func DefaultLangSmithConfig() *LangSmithConfig {
	return &LangSmithConfig{
		Enabled:     os.Getenv("LANGCHAIN_TRACING_V2") == "true",
		APIKey:      os.Getenv("LANGCHAIN_API_KEY"),
		APIUrl:      getEnvOrDefault("LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"),
		ProjectName: getEnvOrDefault("LANGCHAIN_PROJECT", "default"),
		SelfHosted:  false,
	}
}

func LoadConfig(configPath string) (*LangSmithConfig, error) {
	if configPath == "" {
		locations := []string{
			"langsmith.yaml",
			"langsmith.yml",
			"config/langsmith.yaml",
			"config/langsmith.yml",
			filepath.Join(os.Getenv("HOME"), ".langsmith", "config.yaml"),
			filepath.Join(os.Getenv("HOME"), ".langsmith", "config.yml"),
		}

		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				configPath = loc
				break
			}
		}
	}

	if configPath == "" {
		return DefaultLangSmithConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config LangSmithConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if os.Getenv("LANGCHAIN_TRACING_V2") != "" {
		config.Enabled = os.Getenv("LANGCHAIN_TRACING_V2") == "true"
	}
	if os.Getenv("LANGCHAIN_API_KEY") != "" {
		config.APIKey = os.Getenv("LANGCHAIN_API_KEY")
	}
	if os.Getenv("LANGCHAIN_ENDPOINT") != "" {
		config.APIUrl = os.Getenv("LANGCHAIN_ENDPOINT")
	}
	if os.Getenv("LANGCHAIN_PROJECT") != "" {
		config.ProjectName = os.Getenv("LANGCHAIN_PROJECT")
	}
	if os.Getenv("LANGSMITH_SELF_HOSTED") != "" {
		config.SelfHosted = os.Getenv("LANGSMITH_SELF_HOSTED") == "true"
	}

	return &config, nil
}

func SaveConfig(config *LangSmithConfig, configPath string) error {
	if configPath == "" {
		configDir := filepath.Join(os.Getenv("HOME"), ".langsmith")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
		configPath = filepath.Join(configDir, "config.yaml")
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func CreateSelfHostedConfig(url, adminUsername, adminPassword, licenseKey, databaseURL, redisURL string) *SelfHostedConfig {
	return &SelfHostedConfig{
		URL:           url,
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
		LicenseKey:    licenseKey,
		DatabaseURL:   databaseURL,
		RedisURL:      redisURL,
		StorageType:   "filesystem",
		StorageConfig: map[string]interface{}{
			"path": "/data/langsmith",
		},
	}
}

func CreateSelfHostedConfigFromEnv() *SelfHostedConfig {
	return &SelfHostedConfig{
		URL:           getEnvOrDefault("LANGSMITH_URL", "http://langsmith-api:8000"),
		AdminUsername: getEnvOrDefault("LANGSMITH_ADMIN_USERNAME", "admin"),
		AdminPassword: os.Getenv("LANGSMITH_ADMIN_PASSWORD"),
		LicenseKey:    os.Getenv("LANGSMITH_LICENSE_KEY"),
		DatabaseURL:   getEnvOrDefault("LANGSMITH_DATABASE_URL", "postgresql://langsmith:password@langsmith-postgres:5432/langsmith"),
		RedisURL:      getEnvOrDefault("LANGSMITH_REDIS_URL", "redis://langsmith-redis:6379/0"),
		StorageType:   getEnvOrDefault("LANGSMITH_STORAGE_TYPE", "filesystem"),
		StorageConfig: parseStorageConfig(),
	}
}

func parseStorageConfig() map[string]interface{} {
	storageType := getEnvOrDefault("LANGSMITH_STORAGE_TYPE", "filesystem")
	config := make(map[string]interface{})

	switch storageType {
	case "s3":
		config["bucket"] = getEnvOrDefault("LANGSMITH_S3_BUCKET", "langsmith")
		config["region"] = getEnvOrDefault("LANGSMITH_S3_REGION", "us-west-2")
		if os.Getenv("LANGSMITH_S3_ENDPOINT") != "" {
			config["endpoint"] = os.Getenv("LANGSMITH_S3_ENDPOINT")
		}
		if os.Getenv("LANGSMITH_S3_ACCESS_KEY") != "" {
			config["access_key"] = os.Getenv("LANGSMITH_S3_ACCESS_KEY")
		}
		if os.Getenv("LANGSMITH_S3_SECRET_KEY") != "" {
			config["secret_key"] = os.Getenv("LANGSMITH_S3_SECRET_KEY")
		}
	case "azure":
		config["container_name"] = getEnvOrDefault("LANGSMITH_AZURE_CONTAINER", "langsmith")
		if os.Getenv("LANGSMITH_AZURE_ACCOUNT_NAME") != "" {
			config["account_name"] = os.Getenv("LANGSMITH_AZURE_ACCOUNT_NAME")
		}
		if os.Getenv("LANGSMITH_AZURE_ACCOUNT_KEY") != "" {
			config["account_key"] = os.Getenv("LANGSMITH_AZURE_ACCOUNT_KEY")
		}
	case "gcp":
		config["bucket"] = getEnvOrDefault("LANGSMITH_GCP_BUCKET", "langsmith")
		if os.Getenv("LANGSMITH_GCP_PROJECT_ID") != "" {
			config["project_id"] = os.Getenv("LANGSMITH_GCP_PROJECT_ID")
		}
	case "filesystem":
		config["path"] = getEnvOrDefault("LANGSMITH_FILESYSTEM_PATH", "/data/langsmith")
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
