package gitops

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type GuardrailConfig struct {
	EnableLinting        bool
	EnableValidation     bool
	EnableDriftDetection bool
	EnablePolicyChecks   bool
}

type Guardrails struct {
	Config GuardrailConfig
}

func NewGuardrails(config GuardrailConfig) *Guardrails {
	return &Guardrails{
		Config: config,
	}
}

func (g *Guardrails) ValidateTerraform(ctx context.Context, path string) error {
	if !g.Config.EnableValidation {
		return nil
	}

	cmd := exec.CommandContext(ctx, "terraform", "validate", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform validation failed: %s: %w", string(output), err)
	}

	return nil
}

func (g *Guardrails) ValidateKubernetes(ctx context.Context, path string) error {
	if !g.Config.EnableValidation {
		return nil
	}

	cmd := exec.CommandContext(ctx, "kubectl", "apply", "--dry-run=client", "-f", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubernetes validation failed: %s: %w", string(output), err)
	}

	return nil
}

func (g *Guardrails) DetectDrift(ctx context.Context, path string) (bool, error) {
	if !g.Config.EnableDriftDetection {
		return false, nil
	}

	cmd := exec.CommandContext(ctx, "terraform", "plan", "-detailed-exitcode", path)
	err := cmd.Run()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ExitCode() == 2 {
			return true, nil
		}
		return false, fmt.Errorf("drift detection failed: %w", err)
	}

	return false, nil
}

func (g *Guardrails) CheckPolicies(ctx context.Context, path string) error {
	if !g.Config.EnablePolicyChecks {
		return nil
	}

	_, err := exec.LookPath("opa")
	if err != nil {
		return fmt.Errorf("OPA not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, "opa", "eval", "--input", path, "--data", "policy.rego", "data.policy.allow")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("policy check failed: %s: %w", string(output), err)
	}

	if !strings.Contains(string(output), "true") {
		return fmt.Errorf("policy check failed: configuration not allowed")
	}

	return nil
}

func (g *Guardrails) ValidateKubernetesAgainstTerraform(ctx context.Context, k8sPath, tfPath string) error {
	if !g.Config.EnableValidation {
		return nil
	}

	cmd := exec.CommandContext(ctx, "terraform", "show", "-json", tfPath)
	tfOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to extract Kubernetes resources from Terraform: %s: %w", string(tfOutput), err)
	}

	cmd = exec.CommandContext(ctx, "kubectl", "get", "-f", k8sPath, "-o", "json")
	k8sOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to extract Kubernetes resources from manifests: %s: %w", string(k8sOutput), err)
	}

	if len(tfOutput) == 0 || len(k8sOutput) == 0 {
		return fmt.Errorf("failed to compare Kubernetes and Terraform configurations: empty output")
	}

	return nil
}

func (g *Guardrails) LintTerraform(ctx context.Context, path string) error {
	if !g.Config.EnableLinting {
		return nil
	}

	_, err := exec.LookPath("tflint")
	if err != nil {
		return fmt.Errorf("tflint not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, "tflint", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("terraform linting failed: %s: %w", string(output), err)
	}

	return nil
}

func (g *Guardrails) LintKubernetes(ctx context.Context, path string) error {
	if !g.Config.EnableLinting {
		return nil
	}

	_, err := exec.LookPath("kubeval")
	if err != nil {
		return fmt.Errorf("kubeval not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, "kubeval", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubernetes linting failed: %s: %w", string(output), err)
	}

	return nil
}
