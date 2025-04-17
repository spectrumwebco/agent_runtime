package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	
	"github.com/spectrumwebco/agent_runtime/internal/exec"
)

type KataManager struct {
	execer     exec.Execer
	configPath string
	runtimeDir string
}

func NewKataManager(execer exec.Execer, configPath, runtimeDir string) *KataManager {
	if execer == nil {
		execer = exec.DefaultExecer
	}
	return &KataManager{
		execer:     execer,
		configPath: configPath,
		runtimeDir: runtimeDir,
	}
}

func (km *KataManager) CreateContainer(ctx context.Context, name string, options map[string]interface{}) error {
	configPath := filepath.Join(km.runtimeDir, fmt.Sprintf("%s-config.toml", name))
	if err := km.generateKataConfig(configPath, options); err != nil {
		return fmt.Errorf("generate kata config: %w", err)
	}
	
	args := []string{
		"run",
		"--name", name,
		"--runtime", "kata",
		"--runtime-flag", fmt.Sprintf("config=%s", configPath),
		"--detach",
	}
	
	if env, ok := options["env"].(map[string]string); ok {
		for k, v := range env {
			args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	if volumes, ok := options["volumes"].(map[string]string); ok {
		for src, dst := range volumes {
			args = append(args, "--volume", fmt.Sprintf("%s:%s", src, dst))
		}
	}
	
	if ports, ok := options["ports"].(map[string]string); ok {
		for host, container := range ports {
			args = append(args, "--publish", fmt.Sprintf("%s:%s", host, container))
		}
	}
	
	image, ok := options["image"].(string)
	if !ok {
		return fmt.Errorf("image is required")
	}
	args = append(args, image)
	
	if cmd, ok := options["cmd"].([]string); ok {
		args = append(args, cmd...)
	}
	
	command := km.execer.CommandContext(ctx, "podman", args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("create kata container: %w: %s", err, string(output))
	}
	
	return nil
}

func (km *KataManager) StopContainer(ctx context.Context, name string) error {
	command := km.execer.CommandContext(ctx, "podman", "stop", name)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stop kata container: %w: %s", err, string(output))
	}
	
	return nil
}

func (km *KataManager) RemoveContainer(ctx context.Context, name string) error {
	command := km.execer.CommandContext(ctx, "podman", "rm", "-f", name)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remove kata container: %w: %s", err, string(output))
	}
	
	configPath := filepath.Join(km.runtimeDir, fmt.Sprintf("%s-config.toml", name))
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove kata config: %w", err)
	}
	
	return nil
}

func (km *KataManager) ExecInContainer(ctx context.Context, name string, command []string) (string, error) {
	args := append([]string{"exec", name}, command...)
	cmd := km.execer.CommandContext(ctx, "podman", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec in kata container: %w: %s", err, string(output))
	}
	
	return string(output), nil
}

func (km *KataManager) ListContainers(ctx context.Context) ([]string, error) {
	command := km.execer.CommandContext(ctx, "podman", "ps", "--filter", "runtime=kata", "--format", "{{.Names}}")
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("list kata containers: %w: %s", err, string(output))
	}
	
	var containers []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line != "" {
			containers = append(containers, line)
		}
	}
	
	return containers, nil
}

func (km *KataManager) generateKataConfig(configPath string, options map[string]interface{}) error {
	const kataConfigTemplate = `
# Kata Containers configuration for Neovim
[hypervisor.qemu]
path = "/usr/bin/qemu-system-x86_64"
kernel = "/usr/share/kata-containers/vmlinux.container"
initrd = "/usr/share/kata-containers/kata-containers-initrd.img"
machine_type = "q35"
{{if .CPUs}}
default_vcpus = {{.CPUs}}
{{else}}
default_vcpus = 2
{{end}}
{{if .Memory}}
default_memory = {{.Memory}}
{{else}}
default_memory = 2048
{{end}}
disable_block_device_use = false
shared_fs = "virtio-fs"
virtio_fs_daemon = "/usr/libexec/kata-containers/virtiofsd"
valid_hypervisor_paths = [
	"/usr/bin/qemu-system-x86_64"
]

[agent.kata]
{{if .Debug}}
debug = true
{{else}}
debug = false
{{end}}

[runtime]
{{if .Debug}}
enable_debug = true
{{else}}
enable_debug = false
{{end}}
{{if .Trace}}
enable_tracing = true
{{else}}
enable_tracing = false
{{end}}

[netmon]
enable_netmon = false

[factory]
enable_template = true
`

	type kataConfigData struct {
		CPUs   int
		Memory int
		Debug  bool
		Trace  bool
	}
	
	data := kataConfigData{
		CPUs:   2,
		Memory: 2048,
		Debug:  false,
		Trace:  false,
	}
	
	if cpus, ok := options["cpus"].(int); ok {
		data.CPUs = cpus
	}
	
	if memory, ok := options["memory"].(int); ok {
		data.Memory = memory
	}
	
	if debug, ok := options["debug"].(bool); ok {
		data.Debug = debug
	}
	
	if trace, ok := options["trace"].(bool); ok {
		data.Trace = trace
	}
	
	tmpl, err := template.New("kata-config").Parse(kataConfigTemplate)
	if err != nil {
		return fmt.Errorf("parse kata config template: %w", err)
	}
	
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("create kata config file: %w", err)
	}
	defer file.Close()
	
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("execute kata config template: %w", err)
	}
	
	return nil
}
