package k9s

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
	config     *config.Config
	kubeconfig string
	namespace  string
	context    string
	readOnly   bool
	headless   bool
}

type ClientOption func(*Client)

func NewClient(cfg *config.Config, opts ...ClientOption) *Client {
	client := &Client{
		config:   cfg,
		readOnly: false,
		headless: false,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithKubeconfig(kubeconfig string) ClientOption {
	return func(c *Client) {
		c.kubeconfig = kubeconfig
	}
}

func WithNamespace(namespace string) ClientOption {
	return func(c *Client) {
		c.namespace = namespace
	}
}

func WithContext(context string) ClientOption {
	return func(c *Client) {
		c.context = context
	}
}

func WithReadOnly(readOnly bool) ClientOption {
	return func(c *Client) {
		c.readOnly = readOnly
	}
}

func WithHeadless(headless bool) ClientOption {
	return func(c *Client) {
		c.headless = headless
	}
}

func (c *Client) Run(ctx context.Context) error {
	args := []string{}

	if c.kubeconfig != "" {
		args = append(args, "--kubeconfig", c.kubeconfig)
	}

	if c.namespace != "" {
		args = append(args, "-n", c.namespace)
	}

	if c.context != "" {
		args = append(args, "-c", c.context)
	}

	if c.readOnly {
		args = append(args, "--readonly")
	}

	if c.headless {
		args = append(args, "--headless")
	}

	cmd := exec.CommandContext(ctx, "k9s", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (c *Client) RunWithResource(ctx context.Context, resource string) error {
	args := []string{}

	if c.kubeconfig != "" {
		args = append(args, "--kubeconfig", c.kubeconfig)
	}

	if c.namespace != "" {
		args = append(args, "-n", c.namespace)
	}

	if c.context != "" {
		args = append(args, "-c", c.context)
	}

	if c.readOnly {
		args = append(args, "--readonly")
	}

	if c.headless {
		args = append(args, "--headless")
	}

	args = append(args, "-d", resource)

	cmd := exec.CommandContext(ctx, "k9s", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func IsInstalled() bool {
	_, err := exec.LookPath("k9s")
	return err == nil
}

func Install() error {
	if _, err := exec.LookPath("brew"); err == nil {
		cmd := exec.Command("brew", "install", "derailed/k9s/k9s")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if _, err := exec.LookPath("apt"); err == nil {
		cmd := exec.Command("apt", "update")
		cmd.Run()
		cmd = exec.Command("apt", "install", "-y", "k9s")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	if _, err := exec.LookPath("yum"); err == nil {
		cmd := exec.Command("yum", "install", "-y", "k9s")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return fmt.Errorf("could not determine package manager to install k9s")
}

func GetVersion() (string, error) {
	cmd := exec.Command("k9s", "version", "--short")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func GetKubeconfig() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
