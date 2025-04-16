package coder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type TerraformProvider struct {
	config *config.Config
}

func NewTerraformProvider(cfg *config.Config) *TerraformProvider {
	return &TerraformProvider{
		config: cfg,
	}
}

func (p *TerraformProvider) CreateTemplate(ctx context.Context, name, description string, variables map[string]interface{}) (*Template, error) {
	return &Template{
		ID:          "template-" + name,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}, nil
}

func (p *TerraformProvider) ApplyTemplate(ctx context.Context, templateID string, variables map[string]interface{}) error {
	return nil
}

func (p *TerraformProvider) DestroyTemplate(ctx context.Context, templateID string) error {
	return nil
}

func (p *TerraformProvider) GenerateKataContainerTemplate(ctx context.Context, name, description string, config *KataContainerConfig) (string, error) {
	templateDir := filepath.Join(os.TempDir(), "coder-templates", name)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create template directory: %v", err)
	}

	mainTF := fmt.Sprintf(`
terraform {
  required_providers {
    coder = {
      source  = "coder/coder"
      version = "0.6.0"
    }
    kata = {
      source  = "katacontainers/kata"
      version = "0.1.0"
    }
  }
}

provider "coder" {
  url = "%s"
}

provider "kata" {
  url = "%s"
}

resource "kata_container" "dev" {
  name     = "%s"
  image    = "%s"
  cpu      = %d
  memory   = %d
  disk     = %d
}

resource "coder_agent" "dev" {
  os   = "linux"
  arch = "amd64"
}

resource "coder_app" "code-server" {
  agent_id = coder_agent.dev.id
  name     = "code-server"
  url      = "http://localhost:8080"
  icon     = "/icon/code.svg"
}

resource "coder_metadata" "container_info" {
  resource_id = kata_container.dev.id
  item {
    key   = "image"
    value = kata_container.dev.image
  }
  item {
    key   = "cpu"
    value = kata_container.dev.cpu
  }
  item {
    key   = "memory"
    value = kata_container.dev.memory
  }
  item {
    key   = "disk"
    value = kata_container.dev.disk
  }
}
`, p.config.Coder.URL, p.config.Kata.URL, name, config.Image, config.CPU, config.Memory, config.Disk)

	mainTFPath := filepath.Join(templateDir, "main.tf")
	if err := os.WriteFile(mainTFPath, []byte(mainTF), 0644); err != nil {
		return "", fmt.Errorf("failed to write main.tf: %v", err)
	}

	return templateDir, nil
}

type KataContainerConfig struct {
	Image  string `json:"image"`
	CPU    int    `json:"cpu"`
	Memory int    `json:"memory"`
	Disk   int    `json:"disk"`
}
