package gitops

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuardrails(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "guardrails-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tfFile := filepath.Join(tempDir, "main.tf")
	err = os.WriteFile(tfFile, []byte(`
provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_namespace" "test" {
  metadata {
    name = "test-namespace"
  }
}
`), 0644)
	assert.NoError(t, err)

	k8sFile := filepath.Join(tempDir, "deployment.yaml")
	err = os.WriteFile(k8sFile, []byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: test
        image: nginx:latest
`), 0644)
	assert.NoError(t, err)

	guardrails := NewGuardrails(GuardrailConfig{
		EnableLinting:        false,
		EnableValidation:     false,
		EnableDriftDetection: false,
		EnablePolicyChecks:   false,
	})

	err = guardrails.ValidateTerraform(context.Background(), tempDir)
	assert.NoError(t, err)

	err = guardrails.ValidateKubernetes(context.Background(), k8sFile)
	assert.NoError(t, err)

	drift, err := guardrails.DetectDrift(context.Background(), tempDir)
	assert.NoError(t, err)
	assert.False(t, drift)

	err = guardrails.CheckPolicies(context.Background(), tempDir)
	assert.NoError(t, err)

	err = guardrails.ValidateKubernetesAgainstTerraform(context.Background(), k8sFile, tfFile)
	assert.NoError(t, err)

	err = guardrails.LintTerraform(context.Background(), tempDir)
	assert.NoError(t, err)

	err = guardrails.LintKubernetes(context.Background(), k8sFile)
	assert.NoError(t, err)
}

func TestGuardrailsConfig(t *testing.T) {
	config := GuardrailConfig{
		EnableLinting:        true,
		EnableValidation:     true,
		EnableDriftDetection: true,
		EnablePolicyChecks:   true,
	}

	guardrails := NewGuardrails(config)
	assert.NotNil(t, guardrails)
	assert.Equal(t, config, guardrails.Config)
}
