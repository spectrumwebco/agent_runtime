package rex

import (
	"context"
	"fmt"
)

type DeploymentType string

const (
	DeploymentTypeLocal DeploymentType = "local"
	
	DeploymentTypeDocker DeploymentType = "docker"
	
	DeploymentTypeRemote DeploymentType = "remote"
)

type DeploymentConfig struct {
	Type            DeploymentType `json:"type"`
	DockerConfig    *DockerConfig  `json:"docker_config,omitempty"`
	RemoteConfig    *RemoteConfig  `json:"remote_config,omitempty"`
	PythonStandaloneDir string     `json:"python_standalone_dir,omitempty"`
}

type DockerConfig struct {
	Image       string            `json:"image"`
	Volumes     map[string]string `json:"volumes,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Network     string            `json:"network,omitempty"`
}

type RemoteConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	KeyFile  string `json:"key_file,omitempty"`
}

type Deployment interface {
	Start(ctx context.Context) error
	
	Stop(ctx context.Context) error
	
	GetRuntime() Runtime
}

type LocalDeployment struct {
	runtime *LocalRuntime
	logger  Logger
}

func NewLocalDeployment(logger Logger) *LocalDeployment {
	if logger == nil {
		logger = &DefaultLogger{}
	}
	
	return &LocalDeployment{
		runtime: NewLocalRuntime(logger),
		logger:  logger,
	}
}

func (d *LocalDeployment) Start(ctx context.Context) error {
	d.logger.Info("Starting local deployment")
	return nil
}

func (d *LocalDeployment) Stop(ctx context.Context) error {
	d.logger.Info("Stopping local deployment")
	_, err := d.runtime.Close(ctx)
	return err
}

func (d *LocalDeployment) GetRuntime() Runtime {
	return d.runtime
}

func GetDeployment(config *DeploymentConfig, logger Logger) (Deployment, error) {
	switch config.Type {
	case DeploymentTypeLocal:
		return NewLocalDeployment(logger), nil
	case DeploymentTypeDocker:
		return nil, fmt.Errorf("docker deployment not implemented yet")
	case DeploymentTypeRemote:
		return nil, fmt.Errorf("remote deployment not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported deployment type: %s", config.Type)
	}
}
