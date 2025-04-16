package kops

import (
	"context"
	"io"

	"github.com/spectrumwebco/agent_runtime/internal/kubernetes/kops"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config *config.Config
	client *kops.Client
}

func NewModule(cfg *config.Config) (*Module, error) {
	client, err := kops.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Module{
		config: cfg,
		client: client,
	}, nil
}

func (m *Module) Name() string {
	return "kops"
}

func (m *Module) Description() string {
	return "Kubernetes Operations (kOps) for cluster management"
}

func (m *Module) Initialize(ctx context.Context) error {
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) GetClient() *kops.Client {
	return m.client
}

func (m *Module) SetStateStore(stateStore string) {
	m.client.SetStateStore(stateStore)
}

func (m *Module) SetKubeConfig(kubeconfig string) {
	m.client.SetKubeConfig(kubeconfig)
}

func (m *Module) SetClusterName(clusterName string) {
	m.client.SetClusterName(clusterName)
}

func (m *Module) CreateCluster(ctx context.Context, options kops.CreateClusterOptions) error {
	return m.client.CreateCluster(ctx, options)
}

func (m *Module) UpdateCluster(ctx context.Context, options kops.UpdateClusterOptions) error {
	return m.client.UpdateCluster(ctx, options)
}

func (m *Module) DeleteCluster(ctx context.Context, options kops.DeleteClusterOptions) error {
	return m.client.DeleteCluster(ctx, options)
}

func (m *Module) ValidateCluster(ctx context.Context, options kops.ValidateClusterOptions) error {
	return m.client.ValidateCluster(ctx, options)
}

func (m *Module) GetClusters(ctx context.Context, options kops.GetClustersOptions) error {
	return m.client.GetClusters(ctx, options)
}

func (m *Module) ExportKubecfg(ctx context.Context, options kops.ExportKubecfgOptions) error {
	return m.client.ExportKubecfg(ctx, options)
}

func (m *Module) RollingUpdate(ctx context.Context, options kops.RollingUpdateOptions) error {
	return m.client.RollingUpdate(ctx, options)
}

func (m *Module) GetInstanceGroups(ctx context.Context, options kops.GetInstanceGroupsOptions) error {
	return m.client.GetInstanceGroups(ctx, options)
}

func (m *Module) EditInstanceGroup(ctx context.Context, options kops.EditInstanceGroupOptions) error {
	return m.client.EditInstanceGroup(ctx, options)
}

func (m *Module) GetSecrets(ctx context.Context, options kops.GetSecretsOptions) error {
	return m.client.GetSecrets(ctx, options)
}

func (m *Module) ToolboxDump(ctx context.Context, options kops.ToolboxDumpOptions) error {
	return m.client.ToolboxDump(ctx, options)
}

func (m *Module) InstallBinary(ctx context.Context, version string, installPath string) error {
	return m.client.InstallBinary(ctx, version, installPath)
}

func (m *Module) RunExample(ctx context.Context, stdout io.Writer, stderr io.Writer) error {
	m.SetStateStore("s3://kops-state-store")
	
	return m.GetClusters(ctx, kops.GetClustersOptions{
		Output: "json",
		Stdout: stdout,
		Stderr: stderr,
	})
}
