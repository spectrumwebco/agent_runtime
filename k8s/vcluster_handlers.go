package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type VClusterConfig struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	KubernetesVersion string          `json:"kubernetesVersion"`
	Helm            HelmConfig        `json:"helm,omitempty"`
	Values          map[string]interface{} `json:"values,omitempty"`
}

type HelmConfig struct {
	Chart       string `json:"chart,omitempty"`
	ChartVersion string `json:"chartVersion,omitempty"`
	Repository  string `json:"repository,omitempty"`
}

func RegisterVClusterTools(mcpServer *server.MCPServer) {
	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_vclusters",
		Description: "List all vClusters",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listVClusters(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "create_vcluster",
		Description: "Create a new vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "kubernetes_version",
				Description: "Kubernetes version for the vCluster",
				Type:        "string",
				Required:    false,
			},
			{
				Name:        "values",
				Description: "Custom values for the vCluster Helm chart",
				Type:        "object",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return createVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_vcluster",
		Description: "Delete a vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return deleteVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "connect_to_vcluster",
		Description: "Connect to a vCluster and get kubeconfig",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return connectToVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "disconnect_from_vcluster",
		Description: "Disconnect from a vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return disconnectFromVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "upgrade_vcluster",
		Description: "Upgrade a vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "kubernetes_version",
				Description: "Kubernetes version for the vCluster",
				Type:        "string",
				Required:    false,
			},
			{
				Name:        "values",
				Description: "Custom values for the vCluster Helm chart",
				Type:        "object",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return upgradeVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "pause_vcluster",
		Description: "Pause a vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return pauseVCluster(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "resume_vcluster",
		Description: "Resume a paused vCluster",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "vCluster name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return resumeVCluster(ctx, params)
		},
	})
}

func listVClusters(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	args := []string{"list", "-o", "json"}
	
	if params["namespace"] != nil {
		namespace := params["namespace"].(string)
		args = append(args, "-n", namespace)
	}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list vClusters: %w, output: %s", err, string(output))
	}
	
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse vCluster list output: %w", err)
	}
	
	return result, nil
}

func createVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"create", name, "-n", namespace}
	
	if params["kubernetes_version"] != nil {
		kubernetesVersion := params["kubernetes_version"].(string)
		args = append(args, "--kubernetes-version", kubernetesVersion)
	}
	
	if params["values"] != nil {
		valuesMap, ok := params["values"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid values format")
		}
		
		valuesBytes, err := json.Marshal(valuesMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal values: %w", err)
		}
		
		tempDir, err := os.MkdirTemp("", "vcluster-values-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
		
		valuesFile := filepath.Join(tempDir, "values.yaml")
		if err := os.WriteFile(valuesFile, valuesBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to write values file: %w", err)
		}
		
		args = append(args, "-f", valuesFile)
	}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func deleteVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"delete", name, "-n", namespace}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to delete vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func connectToVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	tempDir, err := os.MkdirTemp("", "vcluster-kubeconfig-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	kubeconfigPath := filepath.Join(tempDir, "kubeconfig.yaml")
	
	args := []string{"connect", name, "-n", namespace, "--kube-config", kubeconfigPath, "--update-current=false"}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to connect to vCluster: %w, output: %s", err, string(output))
	}
	
	kubeconfig, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}
	
	os.RemoveAll(tempDir)
	
	return map[string]interface{}{
		"success":    true,
		"name":       name,
		"namespace":  namespace,
		"kubeconfig": string(kubeconfig),
	}, nil
}

func disconnectFromVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"disconnect", name, "-n", namespace}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to disconnect from vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func upgradeVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"upgrade", name, "-n", namespace}
	
	if params["kubernetes_version"] != nil {
		kubernetesVersion := params["kubernetes_version"].(string)
		args = append(args, "--kubernetes-version", kubernetesVersion)
	}
	
	if params["values"] != nil {
		valuesMap, ok := params["values"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid values format")
		}
		
		valuesBytes, err := json.Marshal(valuesMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal values: %w", err)
		}
		
		tempDir, err := os.MkdirTemp("", "vcluster-values-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer os.RemoveAll(tempDir)
		
		valuesFile := filepath.Join(tempDir, "values.yaml")
		if err := os.WriteFile(valuesFile, valuesBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to write values file: %w", err)
		}
		
		args = append(args, "-f", valuesFile)
	}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func pauseVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"pause", name, "-n", namespace}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to pause vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func resumeVCluster(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"resume", name, "-n", namespace}
	
	cmd := exec.CommandContext(ctx, "vcluster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to resume vCluster: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}
