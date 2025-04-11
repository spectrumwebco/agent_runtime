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

type PolicyConfig struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Source    string            `json:"source"`
	Type      string            `json:"type"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func RegisterJSPolicyTools(mcpServer *server.MCPServer) {
	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_policies",
		Description: "List all jsPolicies",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listPolicies(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "get_policy",
		Description: "Get details of a specific jsPolicy",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Policy name",
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
			return getPolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "create_policy",
		Description: "Create a new jsPolicy",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Policy name",
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
				Name:        "source",
				Description: "Policy source code (JavaScript)",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "type",
				Description: "Policy type (ValidatingAdmissionPolicy, MutatingAdmissionPolicy)",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "labels",
				Description: "Labels for the policy",
				Type:        "object",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return createPolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "update_policy",
		Description: "Update an existing jsPolicy",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Policy name",
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
				Name:        "source",
				Description: "Policy source code (JavaScript)",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return updatePolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_policy",
		Description: "Delete a jsPolicy",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Policy name",
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
			return deletePolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "validate_policy",
		Description: "Validate a jsPolicy",
		Parameters: []mcp.Parameter{
			{
				Name:        "source",
				Description: "Policy source code (JavaScript)",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return validatePolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "test_policy",
		Description: "Test a jsPolicy against a resource",
		Parameters: []mcp.Parameter{
			{
				Name:        "source",
				Description: "Policy source code (JavaScript)",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "resource",
				Description: "Resource to test against (JSON)",
				Type:        "object",
				Required:    true,
			},
			{
				Name:        "operation",
				Description: "Operation (CREATE, UPDATE, DELETE)",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return testPolicy(ctx, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "get_policy_violations",
		Description: "Get policy violations",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return getPolicyViolations(ctx, params)
		},
	})
}

func listPolicies(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	args := []string{"policy", "list", "-o", "json"}
	
	if params["namespace"] != nil {
		namespace := params["namespace"].(string)
		args = append(args, "-n", namespace)
	}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list policies: %w, output: %s", err, string(output))
	}
	
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse policy list output: %w", err)
	}
	
	return result, nil
}

func getPolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"policy", "get", name, "-n", namespace, "-o", "json"}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get policy: %w, output: %s", err, string(output))
	}
	
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse policy output: %w", err)
	}
	
	return result, nil
}

func createPolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	source := params["source"].(string)
	policyType := params["type"].(string)
	
	tempDir, err := os.MkdirTemp("", "jspolicy-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	policyFile := filepath.Join(tempDir, "policy.js")
	if err := os.WriteFile(policyFile, []byte(source), 0644); err != nil {
		return nil, fmt.Errorf("failed to write policy file: %w", err)
	}
	
	args := []string{"policy", "create", name, "-n", namespace, "--file", policyFile, "--type", policyType}
	
	if params["labels"] != nil {
		labelsMap, ok := params["labels"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid labels format")
		}
		
		for k, v := range labelsMap {
			args = append(args, "--label", fmt.Sprintf("%s=%s", k, v.(string)))
		}
	}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create policy: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"type":    policyType,
		"message": string(output),
	}, nil
}

func updatePolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	source := params["source"].(string)
	
	tempDir, err := os.MkdirTemp("", "jspolicy-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	policyFile := filepath.Join(tempDir, "policy.js")
	if err := os.WriteFile(policyFile, []byte(source), 0644); err != nil {
		return nil, fmt.Errorf("failed to write policy file: %w", err)
	}
	
	args := []string{"policy", "update", name, "-n", namespace, "--file", policyFile}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to update policy: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func deletePolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	namespace := params["namespace"].(string)
	
	args := []string{"policy", "delete", name, "-n", namespace}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to delete policy: %w, output: %s", err, string(output))
	}
	
	return map[string]interface{}{
		"success": true,
		"name":    name,
		"namespace": namespace,
		"message": string(output),
	}, nil
}

func validatePolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	source := params["source"].(string)
	
	tempDir, err := os.MkdirTemp("", "jspolicy-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	policyFile := filepath.Join(tempDir, "policy.js")
	if err := os.WriteFile(policyFile, []byte(source), 0644); err != nil {
		return nil, fmt.Errorf("failed to write policy file: %w", err)
	}
	
	args := []string{"policy", "validate", "--file", policyFile}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return map[string]interface{}{
			"valid":   false,
			"message": string(output),
			"error":   err.Error(),
		}, nil
	}
	
	return map[string]interface{}{
		"valid":   true,
		"message": string(output),
	}, nil
}

func testPolicy(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	source := params["source"].(string)
	resource := params["resource"].(map[string]interface{})
	operation := params["operation"].(string)
	
	tempDir, err := os.MkdirTemp("", "jspolicy-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	policyFile := filepath.Join(tempDir, "policy.js")
	if err := os.WriteFile(policyFile, []byte(source), 0644); err != nil {
		return nil, fmt.Errorf("failed to write policy file: %w", err)
	}
	
	resourceBytes, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}
	
	resourceFile := filepath.Join(tempDir, "resource.json")
	if err := os.WriteFile(resourceFile, resourceBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write resource file: %w", err)
	}
	
	args := []string{"policy", "test", "--file", policyFile, "--resource", resourceFile, "--operation", operation}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	
	result := map[string]interface{}{
		"output": string(output),
	}
	
	if err != nil {
		result["passed"] = false
		result["error"] = err.Error()
	} else {
		result["passed"] = true
	}
	
	if strings.Contains(string(output), "PASSED") {
		result["passed"] = true
	} else if strings.Contains(string(output), "FAILED") {
		result["passed"] = false
	}
	
	return result, nil
}

func getPolicyViolations(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	args := []string{"policy", "violations", "-o", "json"}
	
	if params["namespace"] != nil {
		namespace := params["namespace"].(string)
		args = append(args, "-n", namespace)
	}
	
	cmd := exec.CommandContext(ctx, "jspolicy", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get policy violations: %w, output: %s", err, string(output))
	}
	
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse policy violations output: %w", err)
	}
	
	return result, nil
}
