package rex

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type LocalRuntime struct {
	sessions map[string]*BashSession
	logger   Logger
}

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type DefaultLogger struct{}

func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}

func (l *DefaultLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func (l *DefaultLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

type BashSession struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewLocalRuntime(logger Logger) *LocalRuntime {
	if logger == nil {
		logger = &DefaultLogger{}
	}
	
	return &LocalRuntime{
		sessions: make(map[string]*BashSession),
		logger:   logger,
	}
}

func (r *LocalRuntime) IsAlive(ctx context.Context, timeout *time.Duration) (*IsAliveResponse, error) {
	return &IsAliveResponse{
		IsAlive: true,
		Message: "Local runtime is always alive",
	}, nil
}

func (r *LocalRuntime) CreateSession(ctx context.Context, request *CreateBashSessionRequest) (*CreateBashSessionResponse, error) {
	if request.SessionType != "bash" {
		return nil, fmt.Errorf("unsupported session type: %s", request.SessionType)
	}
	
	session := request.Session
	if session == "" {
		session = "default"
	}
	
	if _, exists := r.sessions[session]; exists {
		return nil, fmt.Errorf("session already exists: %s", session)
	}
	
	cmd := exec.CommandContext(ctx, "bash", "--login")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start bash session: %w", err)
	}
	
	r.sessions[session] = &BashSession{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
	
	var output strings.Builder
	for _, source := range request.StartupSource {
		if _, err := stdin.Write([]byte(fmt.Sprintf("source %s\n", source))); err != nil {
			r.logger.Error("Failed to source startup file: %s", err)
		}
		
		timeout := time.Duration(request.StartupTimeout * float64(time.Second))
		buf := make([]byte, 4096)
		readCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		
		done := make(chan struct{})
		go func() {
			n, err := stdout.Read(buf)
			if err != nil && err != io.EOF {
				r.logger.Error("Failed to read startup output: %s", err)
			}
			output.Write(buf[:n])
			close(done)
		}()
		
		select {
		case <-readCtx.Done():
			r.logger.Warn("Startup command timed out")
		case <-done:
		}
	}
	
	return &CreateBashSessionResponse{
		Output:      output.String(),
		SessionType: "bash",
	}, nil
}

func (r *LocalRuntime) RunInSession(ctx context.Context, action interface{}) (*BashObservation, error) {
	switch a := action.(type) {
	case *BashAction:
		return r.runBashAction(ctx, a)
	case *BashInterruptAction:
		return r.runBashInterruptAction(ctx, a)
	default:
		return nil, fmt.Errorf("unsupported action type: %T", action)
	}
}

func (r *LocalRuntime) runBashAction(ctx context.Context, action *BashAction) (*BashObservation, error) {
	session := action.Session
	if session == "" {
		session = "default"
	}
	
	s, exists := r.sessions[session]
	if !exists {
		return nil, fmt.Errorf("session does not exist: %s", session)
	}
	
	if _, err := s.stdin.Write([]byte(action.Command + "\n")); err != nil {
		return nil, fmt.Errorf("failed to write command to stdin: %w", err)
	}
	
	var timeout time.Duration
	if action.Timeout != nil {
		timeout = time.Duration(*action.Timeout * float64(time.Second))
	} else {
		timeout = 30 * time.Second // Default timeout
	}
	
	var output strings.Builder
	var exitCode *int
	
	readCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	buf := make([]byte, 4096)
	done := make(chan struct{})
	go func() {
		for {
			n, err := s.stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					r.logger.Error("Failed to read command output: %s", err)
				}
				break
			}
			
			output.Write(buf[:n])
			
			outputStr := output.String()
			for _, expect := range action.Expect {
				if strings.Contains(outputStr, expect) {
					close(done)
					return
				}
			}
			
			if strings.HasSuffix(strings.TrimSpace(outputStr), "$ ") {
				close(done)
				return
			}
		}
	}()
	
	select {
	case <-readCtx.Done():
		if errors.Is(readCtx.Err(), context.DeadlineExceeded) {
			r.logger.Warn("Command timed out: %s", action.Command)
			return &BashObservation{
				Output:        output.String(),
				ExitCode:      nil,
				FailureReason: "timeout",
				SessionType:   "bash",
			}, nil
		}
		return nil, readCtx.Err()
	case <-done:
	}
	
	if action.Check != "ignore" {
		if _, err := s.stdin.Write([]byte("echo $?\n")); err != nil {
			r.logger.Error("Failed to write exit code command: %s", err)
		} else {
			exitCodeBuf := make([]byte, 10)
			n, err := s.stdout.Read(exitCodeBuf)
			if err != nil {
				r.logger.Error("Failed to read exit code: %s", err)
			} else {
				exitCodeStr := strings.TrimSpace(string(exitCodeBuf[:n]))
				var code int
				_, err := fmt.Sscanf(exitCodeStr, "%d", &code)
				if err == nil {
					exitCode = &code
				}
			}
		}
	}
	
	return &BashObservation{
		Output:      output.String(),
		ExitCode:    exitCode,
		SessionType: "bash",
	}, nil
}

func (r *LocalRuntime) runBashInterruptAction(ctx context.Context, action *BashInterruptAction) (*BashObservation, error) {
	session := action.Session
	if session == "" {
		session = "default"
	}
	
	s, exists := r.sessions[session]
	if !exists {
		return nil, fmt.Errorf("session does not exist: %s", session)
	}
	
	if s.cmd.Process != nil {
		if err := s.cmd.Process.Signal(syscall.SIGINT); err != nil {
			return nil, fmt.Errorf("failed to send interrupt signal: %w", err)
		}
	}
	
	timeout := time.Duration(action.Timeout * float64(time.Second))
	var output strings.Builder
	
	readCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	buf := make([]byte, 4096)
	done := make(chan struct{})
	go func() {
		n, err := s.stdout.Read(buf)
		if err != nil && err != io.EOF {
			r.logger.Error("Failed to read interrupt output: %s", err)
		}
		output.Write(buf[:n])
		close(done)
	}()
	
	select {
	case <-readCtx.Done():
		if errors.Is(readCtx.Err(), context.DeadlineExceeded) {
			r.logger.Warn("Interrupt timed out")
			return &BashObservation{
				Output:        output.String(),
				FailureReason: "timeout",
				SessionType:   "bash",
			}, nil
		}
		return nil, readCtx.Err()
	case <-done:
	}
	
	return &BashObservation{
		Output:      output.String(),
		SessionType: "bash",
	}, nil
}

func (r *LocalRuntime) CloseSession(ctx context.Context, request *CloseBashSessionRequest) (*CloseBashSessionResponse, error) {
	session := request.Session
	if session == "" {
		session = "default"
	}
	
	s, exists := r.sessions[session]
	if !exists {
		return nil, fmt.Errorf("session does not exist: %s", session)
	}
	
	if _, err := s.stdin.Write([]byte("exit\n")); err != nil {
		r.logger.Error("Failed to write exit command: %s", err)
	}
	
	if err := s.cmd.Wait(); err != nil {
		r.logger.Error("Failed to wait for process: %s", err)
	}
	
	delete(r.sessions, session)
	
	return &CloseBashSessionResponse{
		SessionType: "bash",
	}, nil
}

func (r *LocalRuntime) Execute(ctx context.Context, command *Command) (*CommandResponse, error) {
	var cmd *exec.Cmd
	
	switch c := command.Command.(type) {
	case string:
		if command.Shell {
			cmd = exec.CommandContext(ctx, "bash", "-c", c)
		} else {
			cmd = exec.CommandContext(ctx, c)
		}
	case []string:
		if len(c) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmd = exec.CommandContext(ctx, c[0], c[1:]...)
	default:
		return nil, fmt.Errorf("unsupported command type: %T", command.Command)
	}
	
	if command.Env != nil {
		cmd.Env = os.Environ()
		for k, v := range command.Env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	
	if command.Cwd != "" {
		cmd.Dir = command.Cwd
	}
	
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	var err error
	if command.Timeout != nil {
		timeout := time.Duration(*command.Timeout * float64(time.Second))
		execCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		
		err = cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("failed to start command: %w", err)
		}
		
		done := make(chan error)
		go func() {
			done <- cmd.Wait()
		}()
		
		select {
		case <-execCtx.Done():
			if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				return &CommandResponse{
					Stdout:   stdout.String(),
					Stderr:   stderr.String(),
					ExitCode: nil,
				}, nil
			}
			return nil, execCtx.Err()
		case err = <-done:
		}
	} else {
		err = cmd.Run()
	}
	
	var exitCode *int
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			exitCode = &code
		}
		
		if command.Check {
			return nil, fmt.Errorf("%s: %w", command.ErrorMsg, err)
		}
	} else {
		code := 0
		exitCode = &code
	}
	
	return &CommandResponse{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

func (r *LocalRuntime) ReadFile(ctx context.Context, request *ReadFileRequest) (*ReadFileResponse, error) {
	content, err := tools.ReadFile(request.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return &ReadFileResponse{
		Content: content,
	}, nil
}

func (r *LocalRuntime) WriteFile(ctx context.Context, request *WriteFileRequest) (*WriteFileResponse, error) {
	err := tools.WriteFile(request.Path, request.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	
	return &WriteFileResponse{}, nil
}

func (r *LocalRuntime) Upload(ctx context.Context, request *UploadRequest) (*UploadResponse, error) {
	content, err := tools.ReadFile(request.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}
	
	targetDir := filepath.Dir(request.TargetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create target directory: %w", err)
	}
	
	err = tools.WriteFile(request.TargetPath, content)
	if err != nil {
		return nil, fmt.Errorf("failed to write target file: %w", err)
	}
	
	return &UploadResponse{}, nil
}

func (r *LocalRuntime) Close(ctx context.Context) (*CloseResponse, error) {
	for session, s := range r.sessions {
		if _, err := s.stdin.Write([]byte("exit\n")); err != nil {
			r.logger.Error("Failed to write exit command to session %s: %s", session, err)
		}
		
		if err := s.cmd.Wait(); err != nil {
			r.logger.Error("Failed to wait for session %s: %s", session, err)
		}
	}
	
	r.sessions = make(map[string]*BashSession)
	
	return &CloseResponse{}, nil
}
