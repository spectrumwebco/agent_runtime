package agent

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spectrumwebco/agent_runtime/internal/agent_framework"
)

type EnvironmentConfig struct {
	Deployment           *rex.DeploymentConfig `json:"deployment"`
	Repo                 *RepoConfig           `json:"repo,omitempty"`
	PostStartupCommands  []string              `json:"post_startup_commands,omitempty"`
	PostStartupCommandTimeout int              `json:"post_startup_command_timeout,omitempty"`
	Name                 string                `json:"name,omitempty"`
}

type RepoConfig struct {
	RepoName   string `json:"repo_name"`
	RepoURL    string `json:"repo_url,omitempty"`
	BaseCommit string `json:"base_commit,omitempty"`
	CloneDepth int    `json:"clone_depth,omitempty"`
}

type Environment struct {
	deployment rex.Deployment
	repo       *Repo
	postStartupCommands []string
	postStartupCommandTimeout int
	hooks      []EnvironmentHook
	name       string
	logger     Logger
}

type EnvironmentHook interface {
	OnInit(env *Environment)
	OnStartDeployment()
	OnCopyRepoStarted(repo *Repo)
	OnEnvironmentStartup()
	OnClose()
}

type Repo struct {
	RepoName   string
	RepoURL    string
	BaseCommit string
	CloneDepth int
}

func NewEnvironment(
	deployment rex.Deployment,
	repo *Repo,
	postStartupCommands []string,
	postStartupCommandTimeout int,
	hooks []EnvironmentHook,
	name string,
	logger Logger,
) *Environment {
	if logger == nil {
		logger = &DefaultLogger{}
	}
	
	env := &Environment{
		deployment: deployment,
		repo:       repo,
		postStartupCommands: postStartupCommands,
		postStartupCommandTimeout: postStartupCommandTimeout,
		hooks:      hooks,
		name:       name,
		logger:     logger,
	}
	
	for _, hook := range hooks {
		hook.OnInit(env)
	}
	
	return env
}

func NewEnvironmentFromConfig(config *EnvironmentConfig, logger Logger) (*Environment, error) {
	if logger == nil {
		logger = &DefaultLogger{}
	}
	
	deployment, err := rex.GetDeployment(config.Deployment, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}
	
	var repo *Repo
	if config.Repo != nil {
		repo = &Repo{
			RepoName:   config.Repo.RepoName,
			RepoURL:    config.Repo.RepoURL,
			BaseCommit: config.Repo.BaseCommit,
			CloneDepth: config.Repo.CloneDepth,
		}
	}
	
	name := config.Name
	if name == "" {
		name = "main"
	}
	
	timeout := config.PostStartupCommandTimeout
	if timeout == 0 {
		timeout = 500
	}
	
	return NewEnvironment(
		deployment,
		repo,
		config.PostStartupCommands,
		timeout,
		nil,
		name,
		logger,
	), nil
}

func (e *Environment) Start(ctx context.Context) error {
	if err := e.initDeployment(ctx); err != nil {
		return fmt.Errorf("failed to initialize deployment: %w", err)
	}
	
	if err := e.Reset(ctx); err != nil {
		return fmt.Errorf("failed to reset environment: %w", err)
	}
	
	for _, command := range e.postStartupCommands {
		if _, err := e.Communicate(ctx, command, e.postStartupCommandTimeout, "raise", "Failed to execute post-startup command"); err != nil {
			return fmt.Errorf("failed to execute post-startup command: %w", err)
		}
	}
	
	return nil
}

func (e *Environment) Reset(ctx context.Context) error {
	if _, err := e.Communicate(ctx, "cd /", 10, "raise", "Failed to change directory"); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	
	if err := e.copyRepo(ctx); err != nil {
		return fmt.Errorf("failed to copy repository: %w", err)
	}
	
	if err := e.resetRepository(ctx); err != nil {
		return fmt.Errorf("failed to reset repository: %w", err)
	}
	
	for _, hook := range e.hooks {
		hook.OnEnvironmentStartup()
	}
	
	return nil
}

func (e *Environment) HardReset(ctx context.Context) error {
	if err := e.Close(ctx); err != nil {
		return fmt.Errorf("failed to close environment: %w", err)
	}
	
	if err := e.Start(ctx); err != nil {
		return fmt.Errorf("failed to start environment: %w", err)
	}
	
	return nil
}

func (e *Environment) Close(ctx context.Context) error {
	e.logger.Info("Beginning environment shutdown...")
	
	if err := e.deployment.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop deployment: %w", err)
	}
	
	for _, hook := range e.hooks {
		hook.OnClose()
	}
	
	return nil
}

func (e *Environment) Communicate(
	ctx context.Context,
	input string,
	timeout int,
	check string,
	errorMsg string,
) (string, error) {
	e.logger.Debug("Input:\n%s", input)
	
	rexCheck := "silent"
	if check == "ignore" {
		rexCheck = "ignore"
	}
	
	timeoutFloat := float64(timeout)
	action := &rex.BashAction{
		Command:    input,
		Timeout:    &timeoutFloat,
		Check:      rexCheck,
		ActionType: "bash",
	}
	
	observation, err := e.deployment.GetRuntime().RunInSession(ctx, action)
	if err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	
	output := observation.Output
	e.logger.Debug("Output:\n%s", output)
	
	if check != "ignore" && observation.ExitCode != nil && *observation.ExitCode != 0 {
		e.logger.Error("%s:\n%s", errorMsg, output)
		msg := fmt.Sprintf("Command %q failed (exit_code=%d): %s", input, *observation.ExitCode, errorMsg)
		e.logger.Error(msg)
		
		if check == "raise" {
			if err := e.Close(ctx); err != nil {
				e.logger.Error("Failed to close environment: %s", err)
			}
			return output, fmt.Errorf(msg)
		}
	}
	
	return output, nil
}

func (e *Environment) ReadFile(ctx context.Context, path string, encoding, errors string) (string, error) {
	request := &rex.ReadFileRequest{
		Path:     path,
		Encoding: encoding,
		Errors:   errors,
	}
	
	response, err := e.deployment.GetRuntime().ReadFile(ctx, request)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	return response.Content, nil
}

func (e *Environment) WriteFile(ctx context.Context, path, content string) error {
	request := &rex.WriteFileRequest{
		Path:    path,
		Content: content,
	}
	
	_, err := e.deployment.GetRuntime().WriteFile(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

func (e *Environment) SetEnvVariables(ctx context.Context, envVariables map[string]string) error {
	if len(envVariables) == 0 {
		e.logger.Debug("No environment variables to set")
		return nil
	}
	
	var envSetters []string
	for k, v := range envVariables {
		envSetters = append(envSetters, fmt.Sprintf("export %s=%q", k, v))
	}
	
	command := strings.Join(envSetters, " && ")
	_, err := e.Communicate(ctx, command, 10, "raise", "Failed to set environment variables")
	if err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}
	
	return nil
}

func (e *Environment) ExecuteCommand(
	ctx context.Context,
	command string,
	shell bool,
	check bool,
	env map[string]string,
	cwd string,
) (*rex.CommandResponse, error) {
	cmd := &rex.Command{
		Command: command,
		Shell:   shell,
		Check:   check,
		Env:     env,
		Cwd:     cwd,
	}
	
	return e.deployment.GetRuntime().Execute(ctx, cmd)
}

func (e *Environment) InterruptSession(ctx context.Context) error {
	e.logger.Info("Interrupting session")
	
	action := &rex.BashInterruptAction{
		Session:    "default",
		Timeout:    0.2,
		NRetry:     3,
		ActionType: "bash_interrupt",
	}
	
	_, err := e.deployment.GetRuntime().RunInSession(ctx, action)
	if err != nil {
		return fmt.Errorf("failed to interrupt session: %w", err)
	}
	
	return nil
}

func (e *Environment) AddHook(hook EnvironmentHook) {
	hook.OnInit(e)
	e.hooks = append(e.hooks, hook)
}

func (e *Environment) initDeployment(ctx context.Context) error {
	for _, hook := range e.hooks {
		hook.OnStartDeployment()
	}
	
	if err := e.deployment.Start(ctx); err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}
	
	request := &rex.CreateBashSessionRequest{
		StartupSource:  []string{"/root/.bashrc"},
		Session:        "default",
		SessionType:    "bash",
		StartupTimeout: 1.0,
	}
	
	_, err := e.deployment.GetRuntime().CreateSession(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	
	if err := e.SetEnvVariables(ctx, map[string]string{
		"LANG":   "C.UTF-8",
		"LC_ALL": "C.UTF-8",
	}); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}
	
	e.logger.Info("Environment Initialized")
	return nil
}

func (e *Environment) copyRepo(ctx context.Context) error {
	if e.repo == nil {
		return nil
	}
	
	output, err := e.Communicate(ctx, "ls", 10, "raise", "Failed to list directory")
	if err != nil {
		return fmt.Errorf("failed to list directory: %w", err)
	}
	
	folders := strings.Split(output, "\n")
	for _, folder := range folders {
		if folder == e.repo.RepoName {
			return nil
		}
	}
	
	for _, hook := range e.hooks {
		hook.OnCopyRepoStarted(e.repo)
	}
	
	if e.repo.RepoURL != "" {
		cloneDepth := ""
		if e.repo.CloneDepth > 0 {
			cloneDepth = fmt.Sprintf("--depth %d", e.repo.CloneDepth)
		}
		
		cloneCmd := fmt.Sprintf("git clone %s %s /%s", cloneDepth, e.repo.RepoURL, e.repo.RepoName)
		_, err := e.Communicate(ctx, cloneCmd, 300, "raise", "Failed to clone repository")
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	
	return nil
}

func (e *Environment) resetRepository(ctx context.Context) error {
	if e.repo == nil {
		return nil
	}
	
	e.logger.Debug("Resetting repository %s to commit %s", e.repo.RepoName, e.repo.BaseCommit)
	
	startupCommands := []string{
		fmt.Sprintf("cd /%s", e.repo.RepoName),
		"export ROOT=$(pwd -P)",
	}
	
	if e.repo.BaseCommit != "" {
		startupCommands = append(startupCommands,
			"git fetch --all",
			"git reset --hard HEAD",
			"git clean -fdx",
			fmt.Sprintf("git checkout %s", e.repo.BaseCommit),
		)
	}
	
	_, err := e.Communicate(
		ctx,
		strings.Join(startupCommands, " && "),
		120,
		"raise",
		"Failed to clean repository",
	)
	if err != nil {
		return fmt.Errorf("failed to reset repository: %w", err)
	}
	
	return nil
}
