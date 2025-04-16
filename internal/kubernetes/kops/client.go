package kops

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Client struct {
	config      *config.Config
	binaryPath  string
	stateStore  string
	kubeconfig  string
	clusterName string
}

func NewClient(cfg *config.Config) (*Client, error) {
	binaryPath, err := exec.LookPath("kops")
	if err != nil {
		return nil, fmt.Errorf("kops binary not found in PATH: %v", err)
	}

	client := &Client{
		config:     cfg,
		binaryPath: binaryPath,
	}

	return client, nil
}

func (c *Client) SetStateStore(stateStore string) {
	c.stateStore = stateStore
}

func (c *Client) SetKubeConfig(kubeconfig string) {
	c.kubeconfig = kubeconfig
}

func (c *Client) SetClusterName(clusterName string) {
	c.clusterName = clusterName
}

func (c *Client) CreateCluster(ctx context.Context, options CreateClusterOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before creating a cluster")
	}

	args := []string{
		"create", "cluster",
		"--name", options.Name,
		"--state", c.stateStore,
		"--zones", strings.Join(options.Zones, ","),
		"--node-count", fmt.Sprintf("%d", options.NodeCount),
	}

	if options.MasterCount > 0 {
		args = append(args, "--master-count", fmt.Sprintf("%d", options.MasterCount))
	}

	if options.NodeSize != "" {
		args = append(args, "--node-size", options.NodeSize)
	}

	if options.MasterSize != "" {
		args = append(args, "--master-size", options.MasterSize)
	}

	if options.CloudProvider != "" {
		args = append(args, "--cloud", options.CloudProvider)
	}

	if options.NetworkCIDR != "" {
		args = append(args, "--network-cidr", options.NetworkCIDR)
	}

	if options.KubernetesVersion != "" {
		args = append(args, "--kubernetes-version", options.KubernetesVersion)
	}

	if options.SSHPublicKey != "" {
		args = append(args, "--ssh-public-key", options.SSHPublicKey)
	}

	if options.DryRun {
		args = append(args, "--dry-run")
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	if options.Yes {
		args = append(args, "--yes")
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) UpdateCluster(ctx context.Context, options UpdateClusterOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before updating a cluster")
	}

	if c.clusterName == "" && options.Name == "" {
		return fmt.Errorf("cluster name must be set before updating a cluster")
	}

	clusterName := c.clusterName
	if options.Name != "" {
		clusterName = options.Name
	}

	args := []string{
		"update", "cluster",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Yes {
		args = append(args, "--yes")
	}

	if options.CreateKubeconfig {
		args = append(args, "--create-kube-config")
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) DeleteCluster(ctx context.Context, options DeleteClusterOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before deleting a cluster")
	}

	if c.clusterName == "" && options.Name == "" {
		return fmt.Errorf("cluster name must be set before deleting a cluster")
	}

	clusterName := c.clusterName
	if options.Name != "" {
		clusterName = options.Name
	}

	args := []string{
		"delete", "cluster",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Yes {
		args = append(args, "--yes")
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) ValidateCluster(ctx context.Context, options ValidateClusterOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before validating a cluster")
	}

	if c.clusterName == "" && options.Name == "" {
		return fmt.Errorf("cluster name must be set before validating a cluster")
	}

	clusterName := c.clusterName
	if options.Name != "" {
		clusterName = options.Name
	}

	args := []string{
		"validate", "cluster",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) GetClusters(ctx context.Context, options GetClustersOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before listing clusters")
	}

	args := []string{
		"get", "clusters",
		"--state", c.stateStore,
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) ExportKubecfg(ctx context.Context, options ExportKubecfgOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before exporting kubeconfig")
	}

	if c.clusterName == "" && options.Name == "" {
		return fmt.Errorf("cluster name must be set before exporting kubeconfig")
	}

	clusterName := c.clusterName
	if options.Name != "" {
		clusterName = options.Name
	}

	args := []string{
		"export", "kubecfg",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.KubeConfig != "" {
		args = append(args, "--kubeconfig", options.KubeConfig)
	}

	if options.Admin {
		args = append(args, "--admin")
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) RollingUpdate(ctx context.Context, options RollingUpdateOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before performing a rolling update")
	}

	if c.clusterName == "" && options.Name == "" {
		return fmt.Errorf("cluster name must be set before performing a rolling update")
	}

	clusterName := c.clusterName
	if options.Name != "" {
		clusterName = options.Name
	}

	args := []string{
		"rolling-update", "cluster",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Yes {
		args = append(args, "--yes")
	}

	if options.Force {
		args = append(args, "--force")
	}

	if options.CloudOnly {
		args = append(args, "--cloud-only")
	}

	if options.MasterInterval != "" {
		args = append(args, "--master-interval", options.MasterInterval)
	}

	if options.NodeInterval != "" {
		args = append(args, "--node-interval", options.NodeInterval)
	}

	if options.InstanceGroupNames != nil && len(options.InstanceGroupNames) > 0 {
		args = append(args, "--instance-group", strings.Join(options.InstanceGroupNames, ","))
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) GetInstanceGroups(ctx context.Context, options GetInstanceGroupsOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before listing instance groups")
	}

	if c.clusterName == "" && options.ClusterName == "" {
		return fmt.Errorf("cluster name must be set before listing instance groups")
	}

	clusterName := c.clusterName
	if options.ClusterName != "" {
		clusterName = options.ClusterName
	}

	args := []string{
		"get", "instancegroups",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) EditInstanceGroup(ctx context.Context, options EditInstanceGroupOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before editing an instance group")
	}

	if c.clusterName == "" && options.ClusterName == "" {
		return fmt.Errorf("cluster name must be set before editing an instance group")
	}

	if options.InstanceGroupName == "" {
		return fmt.Errorf("instance group name must be provided")
	}

	clusterName := c.clusterName
	if options.ClusterName != "" {
		clusterName = options.ClusterName
	}

	args := []string{
		"edit", "instancegroup", options.InstanceGroupName,
		"--name", clusterName,
		"--state", c.stateStore,
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) GetSecrets(ctx context.Context, options GetSecretsOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before listing secrets")
	}

	args := []string{
		"get", "secrets",
		"--state", c.stateStore,
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) ToolboxDump(ctx context.Context, options ToolboxDumpOptions) error {
	if c.stateStore == "" {
		return fmt.Errorf("state store must be set before dumping cluster state")
	}

	if c.clusterName == "" && options.ClusterName == "" {
		return fmt.Errorf("cluster name must be set before dumping cluster state")
	}

	clusterName := c.clusterName
	if options.ClusterName != "" {
		clusterName = options.ClusterName
	}

	args := []string{
		"toolbox", "dump",
		"--name", clusterName,
		"--state", c.stateStore,
	}

	if options.Output != "" {
		args = append(args, "--output", options.Output)
	}

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	return cmd.Run()
}

func (c *Client) InstallBinary(ctx context.Context, version string, installPath string) error {
	if version == "" {
		version = "latest"
	}

	if installPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %v", err)
		}
		installPath = filepath.Join(homeDir, ".local", "bin")
	}

	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %v", err)
	}

	downloadURL := fmt.Sprintf("https://github.com/kubernetes/kops/releases/download/%s/kops-linux-amd64", version)
	if version == "latest" {
		downloadURL = "https://github.com/kubernetes/kops/releases/latest/download/kops-linux-amd64"
	}

	outputPath := filepath.Join(installPath, "kops")
	cmd := exec.CommandContext(ctx, "curl", "-L", downloadURL, "-o", outputPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download kOps binary: %v", err)
	}

	if err := os.Chmod(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to make kOps binary executable: %v", err)
	}

	return nil
}
