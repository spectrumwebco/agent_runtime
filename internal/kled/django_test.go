package kled

import (
	"context"
	"testing"
)

func TestNewDjangoIntegration(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}

	if djangoIntegration == nil {
		t.Fatal("Django integration is nil")
	}

	if djangoIntegration.Framework != framework {
		t.Errorf("Expected framework to be %v, got %v", framework, djangoIntegration.Framework)
	}

	
	if djangoIntegration.Server.Addr != djangoConfig.Address {
		t.Errorf("Expected server address %s, got %s", djangoConfig.Address, djangoIntegration.Server.Addr)
	}
}

func TestDjangoIntegrationName(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}

	if djangoIntegration.Name() != "django" {
		t.Errorf("Expected name to be %s, got %s", "django", djangoIntegration.Name())
	}
}

func TestDjangoIntegrationStartStop(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}


	if err := djangoIntegration.Start(); err != nil {
		t.Fatalf("Failed to start Django integration: %v", err)
	}

	if err := djangoIntegration.Stop(); err != nil {
		t.Fatalf("Failed to stop Django integration: %v", err)
	}
}

func TestDjangoIntegrationRegisterRoute(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}

	err = djangoIntegration.RegisterRoute("/test", "GET", func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
		return "test", nil
	})

	if err != nil {
		t.Fatalf("Failed to register route: %v", err)
	}

	routes := djangoIntegration.Routes()
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}

	route, ok := routes["/test"]
	if !ok {
		t.Fatalf("Route /test not found")
	}

	result, err := route.Handle(context.Background(), map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to execute route handler: %v", err)
	}

	if result != "test" {
		t.Errorf("Expected result to be %s, got %s", "test", result)
	}
}

func TestDjangoIntegrationRegisterModel(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}

	err = djangoIntegration.RegisterModel("TestModel", map[string]string{
		"id":   "int",
		"name": "string",
	})

	if err != nil {
		t.Fatalf("Failed to register model: %v", err)
	}

	models := djangoIntegration.Models()
	if len(models) != 1 {
		t.Errorf("Expected 1 model, got %d", len(models))
	}

	model, ok := models["TestModel"]
	if !ok {
		t.Fatalf("Model TestModel not found")
	}

	fields := model.Fields()
	if len(fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fields))
	}

	if fields["id"] != "int" {
		t.Errorf("Expected field id to be %s, got %s", "int", fields["id"])
	}

	if fields["name"] != "string" {
		t.Errorf("Expected field name to be %s, got %s", "string", fields["name"])
	}
}

func TestDjangoIntegrationSaveLoadState(t *testing.T) {
	config := FrameworkConfig{
		Name:    "Test Framework",
		Version: "1.0.0",
		Debug:   true,
	}

	framework, err := NewFramework(config)
	if err != nil {
		t.Fatalf("Failed to create framework: %v", err)
	}

	djangoConfig := DjangoIntegrationConfig{
		Address: "localhost:8000",
		Routes:  []DjangoRoute{},
		DjangoConfig: DjangoConfig{
			AppDir:       "/tmp/django_test",
			TemplatesDir: "/tmp/django_test/templates",
			StaticDir:    "/tmp/django_test/static",
			Port:         8000,
			Debug:        true,
		},
	}

	djangoIntegration, err := NewDjangoIntegration(framework, djangoConfig)
	if err != nil {
		t.Fatalf("Failed to create Django integration: %v", err)
	}

	djangoIntegration.Context = map[string]interface{}{
		"test": "value",
	}

	state, err := djangoIntegration.SaveState()
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	djangoIntegration.Context = map[string]interface{}{}

	err = djangoIntegration.LoadState(state)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if djangoIntegration.Context["test"] != "value" {
		t.Errorf("Expected context value to be %s, got %s", "value", djangoIntegration.Context["test"])
	}
}
