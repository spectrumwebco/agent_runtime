package containers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/exec"
	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type PodmanCLILister struct {
	execer exec.Execer
}

var _ interfaces.ContainerLister = (*PodmanCLILister)(nil)

func NewPodman(execer exec.Execer) *PodmanCLILister {
	if execer == nil {
		execer = exec.DefaultExecer
	}
	return &PodmanCLILister{
		execer: execer,
	}
}

func (pcl *PodmanCLILister) List(ctx context.Context) (interfaces.ListContainersResponse, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := pcl.execer.CommandContext(ctx, "podman", "ps", "--all", "--quiet", "--no-trunc")
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err := cmd.Run(); err != nil {
		return interfaces.ListContainersResponse{}, fmt.Errorf("run podman ps: %w: %q", 
			 err, strings.TrimSpace(stderrBuf.String()))
	}

	ids := make([]string, 0)
	scanner := bufio.NewScanner(&stdoutBuf)
	for scanner.Scan() {
		tmp := strings.TrimSpace(scanner.Text())
		if tmp == "" {
			continue
		}
		ids = append(ids, tmp)
	}
	if err := scanner.Err(); err != nil {
		return interfaces.ListContainersResponse{}, fmt.Errorf("scan podman ps output: %w", err)
	}

	res := interfaces.ListContainersResponse{
		Containers: make([]interfaces.Container, 0, len(ids)),
		Warnings:   make([]string, 0),
	}
	podmanPsStderr := strings.TrimSpace(stderrBuf.String())
	if podmanPsStderr != "" {
		res.Warnings = append(res.Warnings, podmanPsStderr)
	}
	if len(ids) == 0 {
		return res, nil
	}

	podmanInspectStdout, podmanInspectStderr, err := runPodmanInspect(ctx, pcl.execer, ids...)
	if err != nil {
		return interfaces.ListContainersResponse{}, fmt.Errorf("run podman inspect: %w: %s", 
			 err, podmanInspectStderr)
	}

	if len(podmanInspectStderr) > 0 {
		res.Warnings = append(res.Warnings, string(podmanInspectStderr))
	}

	containers, warns, err := convertPodmanInspect(podmanInspectStdout)
	if err != nil {
		return interfaces.ListContainersResponse{}, fmt.Errorf("convert podman inspect output: %w", err)
	}
	res.Warnings = append(res.Warnings, warns...)
	res.Containers = append(res.Containers, containers...)

	return res, nil
}

func runPodmanInspect(ctx context.Context, execer exec.Execer, ids ...string) ([]byte, string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	args := append([]string{"inspect"}, ids...)
	cmd := execer.CommandContext(ctx, "podman", args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	return stdoutBuf.Bytes(), stderrBuf.String(), err
}

func convertPodmanInspect(data []byte) ([]interfaces.Container, []string, error) {
	var podmanContainers []podmanContainer
	if err := json.Unmarshal(data, &podmanContainers); err != nil {
		return nil, nil, fmt.Errorf("unmarshal podman inspect output: %w", err)
	}

	containers := make([]interfaces.Container, 0, len(podmanContainers))
	warnings := make([]string, 0)

	for _, pc := range podmanContainers {
		container, warns := convertPodmanContainer(pc)
		warnings = append(warnings, warns...)
		containers = append(containers, container)
	}

	return containers, warnings, nil
}

func convertPodmanContainer(pc podmanContainer) (interfaces.Container, []string) {
	warnings := make([]string, 0)

	createdAt, err := time.Parse(time.RFC3339, pc.Created)
	if err != nil {
		warnings = append(warnings, fmt.Sprintf("parse container created time: %v", err))
		createdAt = time.Time{}
	}

	friendlyName := pc.Name
	if strings.HasPrefix(friendlyName, "/") {
		friendlyName = friendlyName[1:]
	}

	ports := make([]interfaces.ContainerPort, 0)
	for _, port := range pc.NetworkSettings.Ports {
		for _, binding := range port.HostBindings {
			hostPort, err := strconv.ParseUint(binding.HostPort, 10, 16)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("parse host port %q: %v", binding.HostPort, err))
				continue
			}

			containerPort, err := strconv.ParseUint(port.ContainerPort, 10, 16)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("parse container port %q: %v", port.ContainerPort, err))
				continue
			}

			ports = append(ports, interfaces.ContainerPort{
				Description:   port.Description,
				ContainerPort: uint16(containerPort),
				HostPort:      uint16(hostPort),
				HostIP:        binding.HostIP,
				Protocol:      port.Protocol,
			})
		}
	}

	volumes := make(map[string]string)
	for _, mount := range pc.Mounts {
		volumes[mount.Destination] = mount.Source
	}

	return interfaces.Container{
		CreatedAt:    createdAt,
		FriendlyName: friendlyName,
		ID:           pc.ID,
		Image:        pc.Config.Image,
		Labels:       pc.Config.Labels,
		Ports:        ports,
		Running:      pc.State.Running,
		Status:       pc.State.Status,
		Volumes:      volumes,
	}, warnings
}

type podmanContainer struct {
	ID      string    `json:"Id"`
	Created string    `json:"Created"`
	Name    string    `json:"Name"`
	Config  pcConfig  `json:"Config"`
	State   pcState   `json:"State"`
	Mounts  []pcMount `json:"Mounts"`
	NetworkSettings struct {
		Ports map[string]pcPort `json:"Ports"`
	} `json:"NetworkSettings"`
}

type pcConfig struct {
	Image  string            `json:"Image"`
	Labels map[string]string `json:"Labels"`
}

type pcState struct {
	Running bool   `json:"Running"`
	Status  string `json:"Status"`
}

type pcMount struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
}

type pcPort struct {
	ContainerPort string       `json:"ContainerPort"`
	Protocol      string       `json:"Protocol"`
	Description   string       `json:"Description"`
	HostBindings  []pcHostPort `json:"HostBindings"`
}

type pcHostPort struct {
	HostIP   string `json:"HostIP"`
	HostPort string `json:"HostPort"`
}

func (pcl *PodmanCLILister) CreateContainer(ctx context.Context, name string, image string, options map[string]interface{}) (string, error) {
	args := []string{"run", "--name", name}

	interactive, _ := options["interactive"].(bool)
	if !interactive {
		args = append(args, "--detach")
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

	if labels, ok := options["labels"].(map[string]string); ok {
		for k, v := range labels {
			args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
		}
	}

	if workdir, ok := options["workdir"].(string); ok {
		args = append(args, "--workdir", workdir)
	}

	if user, ok := options["user"].(string); ok {
		args = append(args, "--user", user)
	}

	if network, ok := options["network"].(string); ok {
		args = append(args, "--network", network)
	}

	if runtime, ok := options["runtime"].(string); ok {
		args = append(args, "--runtime", runtime)
	}

	if runtimeFlags, ok := options["runtime_flags"].([]string); ok {
		for _, flag := range runtimeFlags {
			args = append(args, "--runtime-flag", flag)
		}
	}

	args = append(args, image)

	if cmd, ok := options["cmd"].([]string); ok {
		args = append(args, cmd...)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	command := pcl.execer.CommandContext(ctx, "podman", args...)
	command.Stdout = &stdoutBuf
	command.Stderr = &stderrBuf

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("run podman container: %w: %s", err, stderrBuf.String())
	}

	return strings.TrimSpace(stdoutBuf.String()), nil
}

func (pcl *PodmanCLILister) StopContainer(ctx context.Context, id string) error {
	var stderrBuf bytes.Buffer
	cmd := pcl.execer.CommandContext(ctx, "podman", "stop", id)
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stop container: %w: %s", err, stderrBuf.String())
	}

	return nil
}

func (pcl *PodmanCLILister) RemoveContainer(ctx context.Context, id string) error {
	var stderrBuf bytes.Buffer
	cmd := pcl.execer.CommandContext(ctx, "podman", "rm", "-f", id)
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("remove container: %w: %s", err, stderrBuf.String())
	}

	return nil
}

func (pcl *PodmanCLILister) ExecInContainer(ctx context.Context, id string, command []string) (string, error) {
	args := append([]string{"exec", id}, command...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := pcl.execer.CommandContext(ctx, "podman", args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec in container: %w: %s", err, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}
