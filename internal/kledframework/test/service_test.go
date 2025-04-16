package test

import (
	"testing"

	"github.com/spectrumwebco/agent_runtime/internal/kledframework"
)

func TestNewService(t *testing.T) {
	service := kledframework.NewService(kledframework.ServiceConfig{
		Name:    "test-service",
		Version: "0.1.0",
		Address: ":8080",
	})

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.Name != "test-service" {
		t.Errorf("Expected service name to be 'test-service', got '%s'", service.Name)
	}

	if service.Version != "0.1.0" {
		t.Errorf("Expected service version to be '0.1.0', got '%s'", service.Version)
	}

	if service.Server.Addr != ":8080" {
		t.Errorf("Expected service address to be ':8080', got '%s'", service.Server.Addr)
	}
}

func TestNewRegistry(t *testing.T) {
	registry := kledframework.NewRegistry()

	if registry == nil {
		t.Fatal("Expected registry to be created, got nil")
	}

	services := registry.ListServices()
	if len(services) != 0 {
		t.Errorf("Expected registry to be empty, got %d services", len(services))
	}
}
