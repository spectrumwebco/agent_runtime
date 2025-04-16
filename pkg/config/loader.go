package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ConfigLoader struct {
	viper      *viper.Viper
	configPath string
	configName string
	configType string
	env        string
}

func NewConfigLoader(configPath, configName, configType, env string) *ConfigLoader {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)

	v.SetEnvPrefix("AGENT")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &ConfigLoader{
		viper:      v,
		configPath: configPath,
		configName: configName,
		configType: configType,
		env:        env,
	}
}

func (l *ConfigLoader) LoadConfig() (*Config, error) {
	if err := l.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	if l.env != "" {
		envConfigName := fmt.Sprintf("%s.%s", l.configName, l.env)
		l.viper.SetConfigName(envConfigName)
		if err := l.viper.MergeInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to merge environment config: %v", err)
			}
		}

		if err := l.viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal environment config: %v", err)
		}
	}

	return &config, nil
}

func (l *ConfigLoader) SaveConfig(config *Config) error {
	var data []byte
	var err error

	switch l.configType {
	case "json":
		data, err = json.MarshalIndent(config, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(config)
	default:
		return fmt.Errorf("unsupported config type: %s", l.configType)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	configFile := filepath.Join(l.configPath, fmt.Sprintf("%s.%s", l.configName, l.configType))
	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func (l *ConfigLoader) GetString(key string) string {
	return l.viper.GetString(key)
}

func (l *ConfigLoader) GetInt(key string) int {
	return l.viper.GetInt(key)
}

func (l *ConfigLoader) GetBool(key string) bool {
	return l.viper.GetBool(key)
}

func (l *ConfigLoader) GetStringSlice(key string) []string {
	return l.viper.GetStringSlice(key)
}

func (l *ConfigLoader) GetStringMap(key string) map[string]interface{} {
	return l.viper.GetStringMap(key)
}

func (l *ConfigLoader) Set(key string, value interface{}) {
	l.viper.Set(key, value)
}

func (l *ConfigLoader) IsSet(key string) bool {
	return l.viper.IsSet(key)
}

func (l *ConfigLoader) GetConfigPath() string {
	return l.configPath
}

func (l *ConfigLoader) GetConfigName() string {
	return l.configName
}

func (l *ConfigLoader) GetConfigType() string {
	return l.configType
}

func (l *ConfigLoader) GetEnv() string {
	return l.env
}

func (l *ConfigLoader) SetEnv(env string) {
	l.env = env
}

func (l *ConfigLoader) LoadFromEnv() (*Config, error) {
	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from env: %v", err)
	}

	return &config, nil
}

func (l *ConfigLoader) LoadFromFile(filePath string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	switch filepath.Ext(filePath) {
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON config: %v", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal YAML config: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", filepath.Ext(filePath))
	}

	return &config, nil
}
