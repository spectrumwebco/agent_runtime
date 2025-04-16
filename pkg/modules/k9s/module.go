package k9s

import (
	"context"
	"fmt"

	"github.com/spectrumwebco/agent_runtime/internal/kubernetes/k9s"
	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Module struct {
	config *config.Config
	client *k9s.Client
}

func NewModule(cfg *config.Config) *Module {
	return &Module{
		config: cfg,
		client: k9s.NewClient(cfg),
	}
}

func (m *Module) Name() string {
	return "k9s"
}

func (m *Module) Description() string {
	return "Terminal UI for Kubernetes cluster management"
}

func (m *Module) Initialize(ctx context.Context) error {
	if !k9s.IsInstalled() {
		return fmt.Errorf("k9s is not installed, please install it using 'kled k9s install'")
	}
	return nil
}

func (m *Module) Shutdown(ctx context.Context) error {
	return nil
}

func (m *Module) GetClient() *k9s.Client {
	return m.client
}

func (m *Module) Run(ctx context.Context) error {
	return m.client.Run(ctx)
}

func (m *Module) RunWithResource(ctx context.Context, resource string) error {
	return m.client.RunWithResource(ctx, resource)
}

func (m *Module) Install() error {
	return k9s.Install()
}

func (m *Module) GetVersion() (string, error) {
	return k9s.GetVersion()
}

func (m *Module) WithKubeconfig(kubeconfig string) *Module {
	m.client = k9s.NewClient(m.config, k9s.WithKubeconfig(kubeconfig))
	return m
}

func (m *Module) WithNamespace(namespace string) *Module {
	m.client = k9s.NewClient(m.config, k9s.WithNamespace(namespace))
	return m
}

func (m *Module) WithContext(context string) *Module {
	m.client = k9s.NewClient(m.config, k9s.WithContext(context))
	return m
}

func (m *Module) WithReadOnly(readOnly bool) *Module {
	m.client = k9s.NewClient(m.config, k9s.WithReadOnly(readOnly))
	return m
}

func (m *Module) WithHeadless(headless bool) *Module {
	m.client = k9s.NewClient(m.config, k9s.WithHeadless(headless))
	return m
}
