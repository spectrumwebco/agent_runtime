package interfaces

import (
	"context"
	"time"
)

type Container struct {
	CreatedAt    time.Time         `json:"created_at"`
	FriendlyName string            `json:"friendly_name"`
	ID           string            `json:"id"`
	Image        string            `json:"image"`
	Labels       map[string]string `json:"labels"`
	Ports        []ContainerPort   `json:"ports"`
	Running      bool              `json:"running"`
	Status       string            `json:"status"`
	Volumes      map[string]string `json:"volumes"`
}

type ContainerPort struct {
	Description   string `json:"description"`
	ContainerPort uint16 `json:"container_port"`
	HostPort      uint16 `json:"host_port"`
	HostIP        string `json:"host_ip"`
	Protocol      string `json:"protocol"`
}

type ListContainersResponse struct {
	Containers []Container `json:"containers"`
	Warnings   []string    `json:"warnings"`
}

type ContainerLister interface {
	List(ctx context.Context) (ListContainersResponse, error)
	
	CreateContainer(ctx context.Context, name string, image string, options map[string]interface{}) (string, error)
	
	StopContainer(ctx context.Context, id string) error
	
	RemoveContainer(ctx context.Context, id string) error
	
	ExecInContainer(ctx context.Context, id string, command []string) (string, error)
}
